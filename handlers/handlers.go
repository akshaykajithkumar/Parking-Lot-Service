package handlers

import (
	"fmt"
	"log"
	"main/models"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

var Db *gorm.DB
var serverInstance *Server

func SetDB(database *gorm.DB) {
	Db = database
}

func InitServer(server *Server) {
	serverInstance = server
}

type Server struct {
	ParkingLots map[string]*models.ParkingLot
	Clients     map[*websocket.Conn]bool
	Mutex       *sync.Mutex
	DB          *gorm.DB
}

// @Summary Create Parking Lot
// @Description Create a new parking lot [first_rate is the rate of first hours(eg : 1st day = 24 hours ) and after_rate is the rate of after hours]
// @Tags parkinglots system
// @Accept json
// @Produce json
// @Param parkinglot body models.SwaggerParkingLot true "Parking Lot"
// @Success 201 {object} models.SwaggerParkingLot
// @Failure 400 {object} models.ErrorResponse
// @Router /parking-lots [post]
func CreateParkingLot(c echo.Context) error {
	var parkingLot models.ParkingLot
	if err := c.Bind(&parkingLot); err != nil {
		c.Logger().Error("Failed to bind request data: ", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error(), Error: nil})
	}

	parkingLot.ID = 0
	for i := range parkingLot.Tariffs {
		parkingLot.Tariffs[i].ID = 0
		parkingLot.Tariffs[i].ParkingLotID = 0
		for j := range parkingLot.Tariffs[i].RatePlans {
			parkingLot.Tariffs[i].RatePlans[j].ID = 0
			parkingLot.Tariffs[i].RatePlans[j].TariffID = 0
		}
	}

	if err := Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&parkingLot).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.Logger().Error("Failed to create parking lot: ", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Error creating parking lot", Error: err.Error()})
	}

	return c.JSON(http.StatusOK, models.ErrorResponse{Message: "successfully created"})
}

// GetParkingLotDetails godoc
// @Summary Get details of a parking lot
// @Description Get detailed information about a specific parking lot
// @Tags parkinglots system
// @Produce json
// @Param id path int true "Parking Lot ID"
// @Success 200 {object} models.SwaggerParkingLot
// @Failure 404 {object} models.ErrorResponse
// @Router /parking-lots/{id} [get]
func GetParkingLotDetails(c echo.Context) error {
	id := c.Param("id")
	var parkingLot models.ParkingLot
	if err := Db.Preload("Tariffs").Preload("Tariffs.RatePlans").First(&parkingLot, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Error getting parking lots", Error: err.Error()})
	}

	simplifiedLot := map[string]interface{}{
		"ID":                   parkingLot.ID,
		"name":                 parkingLot.Name,
		"motorcycle_spots":     parkingLot.MotorcycleSpots,
		"car_spots":            parkingLot.CarSpots,
		"bus_spots":            parkingLot.BusSpots,
		"occupied_motorcycles": parkingLot.OccupiedMotorcycles,
		"occupied_cars":        parkingLot.OccupiedCars,
		"occupied_buses":       parkingLot.OccupiedBuses,
		"tariffs":              make([]map[string]interface{}, 0),
	}

	for _, tariff := range parkingLot.Tariffs {
		simplifiedTariff := map[string]interface{}{
			"ID":             tariff.ID,
			"parking_lot_id": tariff.ParkingLotID,
			"vehicle_type":   tariff.VehicleType,
			"rate_plans":     make([]map[string]interface{}, 0),
		}

		for _, ratePlan := range tariff.RatePlans {
			simplifiedRatePlan := map[string]interface{}{
				"tariff_id":   ratePlan.TariffID,
				"first_hours": ratePlan.FirstHours,
				"first_rate":  ratePlan.FirstRate,
				"after_rate":  ratePlan.AfterRate,
			}
			simplifiedTariff["rate_plans"] = append(simplifiedTariff["rate_plans"].([]map[string]interface{}), simplifiedRatePlan)
		}

		simplifiedLot["tariffs"] = append(simplifiedLot["tariffs"].([]map[string]interface{}), simplifiedTariff)
	}

	return c.JSON(http.StatusOK, simplifiedLot)
}

// ListParkingLots godoc
// @Summary List all parking lots
// @Description Get a list of all parking lots
// @Tags parkinglots system
// @Produce json
// @Success 200 {array} models.ParkingLotSummary
// @Failure 500 {object} models.ErrorResponse
// @Router /parking-lots [get]
func ListParkingLots(c echo.Context) error {
	var parkingLots []models.ParkingLot
	if err := Db.Find(&parkingLots).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Error listing parking lots", Error: err.Error()})
	}

	var simplifiedParkingLots []map[string]interface{}

	for _, lot := range parkingLots {
		simplifiedLot := map[string]interface{}{
			"ID":   lot.ID,
			"name": lot.Name,
		}

		simplifiedParkingLots = append(simplifiedParkingLots, simplifiedLot)
	}

	return c.JSON(http.StatusOK, simplifiedParkingLots)
}

// ParkVehicle godoc
// @Summary Park a vehicle
// @Description Park a vehicle in a specific parking lot
// @Tags parkinglots system
// @Accept json
// @Produce json
// @Param id path int true "Parking Lot ID"
// @Param vehicle body models.VehicleData true "Vehicle Data"
// @Success 200 {object} models.SwaggerTicket
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /parking-lots/{id}/park [post]
func ParkVehicle(c echo.Context) error {
	id := c.Param("id")

	var parkingLot models.ParkingLot
	if err := Db.First(&parkingLot, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Parking lot not found", Error: err.Error()})
	}

	var vehicleData struct {
		VehicleType   string `json:"vehicle_type"`
		VehicleNumber string `json:"vehicle_number"`
	}
	if err := c.Bind(&vehicleData); err != nil {
		c.Logger().Error("Error binding request body: ", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid request body", Error: err.Error()})
	}

	// Validating form values
	if vehicleData.VehicleType == "" || vehicleData.VehicleNumber == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Vehicle type and number are required"})
	}

	entryTime := time.Now()

	// Creating a ticket
	var ticket models.Ticket
	ticket.ParkingLotID = parkingLot.ID
	ticket.VehicleType = vehicleData.VehicleType
	ticket.VehicleNumber = vehicleData.VehicleNumber
	ticket.EntryTime = entryTime

	// updatingg parking lot availability
	switch vehicleData.VehicleType {
	case "motorcycle":
		if parkingLot.OccupiedMotorcycles < parkingLot.MotorcycleSpots {
			parkingLot.OccupiedMotorcycles++
		} else {
			return c.JSON(http.StatusConflict, models.ErrorResponse{Message: "No available motorcycle spots"})
		}
	case "car":
		if parkingLot.OccupiedCars < parkingLot.CarSpots {
			parkingLot.OccupiedCars++
		} else {
			return c.JSON(http.StatusConflict, models.ErrorResponse{Message: "No available car spots"})
		}
	case "bus":
		if parkingLot.OccupiedBuses < parkingLot.BusSpots {
			parkingLot.OccupiedBuses++
		} else {
			return c.JSON(http.StatusConflict, models.ErrorResponse{Message: "No available bus spots"})
		}
	default:
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid vehicle type"})
	}

	if err := Db.Save(&parkingLot).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Error updating parking lot"})
	}

	if err := Db.Create(&ticket).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Error creating ticket"})
	}

	response := map[string]interface{}{
		"ID":               ticket.ID,
		"vehicle_type":     ticket.VehicleType,
		"vehicle_number":   ticket.VehicleNumber,
		"parking_lot_name": parkingLot.Name,
		"entry_time":       ticket.EntryTime,
	}

	// Broadcastingg update to notify others of the change
	availableSpots := map[string]int{
		"occupied_cars":        parkingLot.OccupiedCars,
		"occupied_motorcycles": parkingLot.OccupiedMotorcycles,
		"occupied_buses":       parkingLot.OccupiedBuses,
	}
	BroadcastUpdate(availableSpots)

	return c.JSON(http.StatusOK, response)
}

// UnparkVehicle godoc
// @Summary Unpark a vehicle
// @Description Unpark a vehicle from a specific parking lot
// @Tags parkinglots system
// @Accept json
// @Produce json
// @Param id path int true "Parking Lot ID"
// @Param vehicle body models.UnparkVehicleData true "Vehicle Data"
// @Success 200 {object} models.UnparkResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /parking-lots/{id}/unpark [post]
func UnparkVehicle(c echo.Context) error {
	id := c.Param("id")

	// Finding the parking lot by ID
	var parkingLot models.ParkingLot
	if err := Db.Preload("Tariffs.RatePlans").First(&parkingLot, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Parking lot not found"})
	}

	var req struct {
		VehicleNumber string `json:"vehicle_number"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid request body"})
	}

	// Checking if vehicle number is provided
	if req.VehicleNumber == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Vehicle number is required"})
	}

	// Here finding the active ticket for the specified vehicle in the parking lot
	var ticket models.Ticket
	if err := Db.Where("parking_lot_id = ? AND vehicle_number = ? AND exit_time IS NULL", parkingLot.ID, req.VehicleNumber).First(&ticket).Error; err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Active ticket not found"})
	}

	// Recording the exit time for the ticket
	exitTime := time.Now()
	ticket.ExitTime = &exitTime

	// Calculating duration in hours and round up to the nearest hour
	duration := exitTime.Sub(ticket.EntryTime).Hours()
	fmt.Printf("Parking duration: %f hours\n", duration)

	// Calculating fee based on the parking lot's tariff for the vehicle type
	fee := calculateFee(duration, parkingLot.Tariffs, ticket.VehicleType)

	if err := Db.Save(&ticket).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Error updating ticket"})
	}

	// Decreaseing the occupied count based on vehicle type
	switch ticket.VehicleType {
	case "motorcycle":
		if parkingLot.OccupiedMotorcycles > 0 {
			parkingLot.OccupiedMotorcycles--
		} else {
			return c.JSON(http.StatusConflict, models.ErrorResponse{Message: "No motorcycles to unpark"})
		}
	case "car":
		if parkingLot.OccupiedCars > 0 {
			parkingLot.OccupiedCars--
		} else {
			return c.JSON(http.StatusConflict, models.ErrorResponse{Message: "No cars to unpark"})
		}
	case "bus":
		if parkingLot.OccupiedBuses > 0 {
			parkingLot.OccupiedBuses--
		} else {
			return c.JSON(http.StatusConflict, models.ErrorResponse{Message: "No buses to unpark"})
		}
	default:
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid vehicle type"})
	}

	if err := Db.Save(&parkingLot).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Error updating parking lot"})
	}

	availableSpots := map[string]int{
		"occupied_cars":        parkingLot.OccupiedCars,
		"occupied_motorcycles": parkingLot.OccupiedMotorcycles,
		"occupied_buses":       parkingLot.OccupiedBuses,
	}

	BroadcastUpdate(availableSpots)

	receipt := map[string]interface{}{
		"receipt": map[string]interface{}{
			"ID":             ticket.ID,
			"vehicle_type":   ticket.VehicleType,
			"vehicle_number": ticket.VehicleNumber,
			"parking_lot_id": ticket.ParkingLotID,
			"entry_time":     ticket.EntryTime,
			"exit_time":      ticket.ExitTime,
		},
		"Total Fee": fee,
	}

	return c.JSON(http.StatusOK, receipt)
}

func calculateFee(duration float64, tariffs []models.Tariff, vehicleType string) float64 {

	var tariff models.Tariff
	found := false
	for _, t := range tariffs {
		if t.VehicleType == vehicleType {
			tariff = t
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Tariff for vehicle type %s not found\n", vehicleType)
		return 0
	}

	// Calculating full hours even for fractions (e.g., 0.1 hours should be 1 hour)
	fullHours := int(math.Ceil(duration))

	// Calculating fee based on the rate plans for the vehicle type
	var totalFee float64
	for _, ratePlan := range tariff.RatePlans {
		fmt.Printf("Applying rate plan: %+v for duration: %f\n", ratePlan, duration)
		if fullHours <= ratePlan.FirstHours {
			totalFee = ratePlan.FirstRate
			break
		} else {
			totalFee = ratePlan.FirstRate + float64(fullHours-ratePlan.FirstHours)*ratePlan.AfterRate
		}
	}

	fmt.Printf("Calculated fee: %f for vehicle type: %s, duration: %f\n", totalFee, vehicleType, duration)
	return totalFee
}

// GetAvailableSpots godoc
// @Summary Get available spots
// @Description Get available spots in a specific parking lot
// @Tags parkinglots system
// @Produce json
// @Param id path int true "Parking Lot ID"
// @Success 200 {object} models.AvailableSpotsResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /parking-lots/{id}/available-spots [get]
func GetAvailableSpots(c echo.Context) error {
	id := c.Param("id")
	var parkingLot models.ParkingLot
	if err := Db.First(&parkingLot, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Parking lot not found"})
	}

	availableSpots := map[string]int{
		"motorcycle": parkingLot.MotorcycleSpots - parkingLot.OccupiedMotorcycles,
		"car":        parkingLot.CarSpots - parkingLot.OccupiedCars,
		"bus":        parkingLot.BusSpots - parkingLot.OccupiedBuses,
	}

	BroadcastUpdate(availableSpots)
	return c.JSON(http.StatusOK, availableSpots)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return err
	}
	defer func() {
		ws.Close()
		log.Println("WebSocket closed")
	}()

	serverInstance.Mutex.Lock()
	serverInstance.Clients[ws] = true
	serverInstance.Mutex.Unlock()

	go func() {
		for {
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Ping error:", err)
				serverInstance.Mutex.Lock()
				delete(serverInstance.Clients, ws)
				serverInstance.Mutex.Unlock()
				return
			}
			time.Sleep(30 * time.Second)
		}
	}()

	for {
		var msg string
		if err := ws.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Read error: %v", err)
			}
			serverInstance.Mutex.Lock()
			delete(serverInstance.Clients, ws)
			serverInstance.Mutex.Unlock()
			break
		}
	}

	return nil
}

func BroadcastUpdate(data map[string]int) {
	serverInstance.Mutex.Lock()
	defer serverInstance.Mutex.Unlock()

	for client := range serverInstance.Clients {
		if err := client.WriteJSON(data); err != nil {
			log.Println("WriteJSON error:", err)
			client.Close()
			delete(serverInstance.Clients, client)
		}
	}
}
