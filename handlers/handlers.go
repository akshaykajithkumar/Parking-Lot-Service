package handlers

import (
	"main/models"
	"net/http"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
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
}

func CreateParkingLot(c echo.Context) error {
	var parkingLot models.ParkingLot
	if err := c.Bind(&parkingLot); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}
	if err := Db.Create(&parkingLot).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating parking lot"})
	}
	return c.JSON(http.StatusOK, parkingLot)
}

func GetParkingLotDetails(c echo.Context) error {
	id := c.Param("id")
	var parkingLot models.ParkingLot
	if err := Db.Preload("Tariffs").First(&parkingLot, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Parking lot not found"})
	}
	return c.JSON(http.StatusOK, parkingLot)
}

func ListParkingLots(c echo.Context) error {
	var parkingLots []models.ParkingLot
	if err := Db.Preload("Tariffs").Find(&parkingLots).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error listing parking lots"})
	}
	return c.JSON(http.StatusOK, parkingLots)
}

func ParkVehicle(c echo.Context) error {
	id := c.Param("id")
	var parkingLot models.ParkingLot
	if err := Db.First(&parkingLot, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Parking lot not found"})
	}

	vehicleType := c.FormValue("vehicle_type")
	vehicleNumber := c.FormValue("vehicle_number")
	entryTime := time.Now()

	var ticket models.Ticket
	ticket.ParkingLotID = parkingLot.ID
	ticket.VehicleType = vehicleType
	ticket.VehicleNumber = vehicleNumber
	ticket.EntryTime = entryTime

	switch vehicleType {
	case "motorcycle":
		if parkingLot.OccupiedMotorcycles < parkingLot.MotorcycleSpots {
			parkingLot.OccupiedMotorcycles++
		} else {
			return c.JSON(http.StatusConflict, map[string]string{"message": "No available motorcycle spots"})
		}
	case "car":
		if parkingLot.OccupiedCars < parkingLot.CarSpots {
			parkingLot.OccupiedCars++
		} else {
			return c.JSON(http.StatusConflict, map[string]string{"message": "No available car spots"})
		}
	case "bus":
		if parkingLot.OccupiedBuses < parkingLot.BusSpots {
			parkingLot.OccupiedBuses++
		} else {
			return c.JSON(http.StatusConflict, map[string]string{"message": "No available bus spots"})
		}
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid vehicle type"})
	}

	if err := Db.Save(&parkingLot).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error parking vehicle"})
	}

	if err := Db.Create(&ticket).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating ticket"})
	}

	BroadcastUpdate(&parkingLot)
	return c.JSON(http.StatusOK, ticket)
}

func UnparkVehicle(c echo.Context) error {
	id := c.Param("id")
	var parkingLot models.ParkingLot
	if err := Db.First(&parkingLot, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Parking lot not found"})
	}

	vehicleNumber := c.FormValue("vehicle_number")
	var ticket models.Ticket
	if err := Db.Where("parking_lot_id = ? AND vehicle_number = ? AND exit_time IS NULL", parkingLot.ID, vehicleNumber).First(&ticket).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Active ticket not found"})
	}

	exitTime := time.Now()
	ticket.ExitTime = &exitTime

	duration := exitTime.Sub(ticket.EntryTime).Hours()
	var fee float64
	for _, tariff := range parkingLot.Tariffs {
		if tariff.VehicleType == ticket.VehicleType {
			fee = calculateFee(duration, tariff)
			break
		}
	}

	switch ticket.VehicleType {
	case "motorcycle":
		if parkingLot.OccupiedMotorcycles > 0 {
			parkingLot.OccupiedMotorcycles--
		} else {
			return c.JSON(http.StatusConflict, map[string]string{"message": "No motorcycles to unpark"})
		}
	case "car":
		if parkingLot.OccupiedCars > 0 {
			parkingLot.OccupiedCars--
		} else {
			return c.JSON(http.StatusConflict, map[string]string{"message": "No cars to unpark"})
		}
	case "bus":
		if parkingLot.OccupiedBuses > 0 {
			parkingLot.OccupiedBuses--
		} else {
			return c.JSON(http.StatusConflict, map[string]string{"message": "No buses to unpark"})
		}
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid vehicle type"})
	}

	if err := Db.Save(&parkingLot).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error updating parking lot"})
	}

	if err := Db.Save(&ticket).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error updating ticket"})
	}

	BroadcastUpdate(&parkingLot)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"ticket": ticket,
		"fee":    fee,
	})
}

func calculateFee(duration float64, tariff models.Tariff) float64 {
	// Round up to the nearest hour
	hours := int64(duration + 0.9999)
	if hours <= 24 {
		return float64(hours) * tariff.RatePerHour
	}
	days := hours / 24
	remainingHours := hours % 24
	return float64(days)*tariff.MaxDailyRate + float64(remainingHours)*tariff.RatePerHour
}
func GetAvailableSpots(c echo.Context) error {
	id := c.Param("id")
	var parkingLot models.ParkingLot
	if err := Db.First(&parkingLot, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Parking lot not found"})
	}

	availableSpots := map[string]int{
		"motorcycle": parkingLot.MotorcycleSpots - parkingLot.OccupiedMotorcycles,
		"car":        parkingLot.CarSpots - parkingLot.OccupiedCars,
		"bus":        parkingLot.BusSpots - parkingLot.OccupiedBuses,
	}

	return c.JSON(http.StatusOK, availableSpots)
}

func BroadcastUpdate(parkingLot *models.ParkingLot) {
	availableSpots := map[string]int{
		"motorcycle": parkingLot.MotorcycleSpots - parkingLot.OccupiedMotorcycles,
		"car":        parkingLot.CarSpots - parkingLot.OccupiedCars,
		"bus":        parkingLot.BusSpots - parkingLot.OccupiedBuses,
	}

	serverInstance.Mutex.Lock()
	defer serverInstance.Mutex.Unlock()

	for client := range serverInstance.Clients {
		if err := websocket.JSON.Send(client, availableSpots); err != nil {
			client.Close()
			delete(serverInstance.Clients, client)
		}
	}
}

func HandleWebSocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		serverInstance.Mutex.Lock()
		serverInstance.Clients[ws] = true
		serverInstance.Mutex.Unlock()

		defer func() {
			serverInstance.Mutex.Lock()
			delete(serverInstance.Clients, ws)
			serverInstance.Mutex.Unlock()
			ws.Close()
		}()

		for {
			var msg string
			if err := websocket.Message.Receive(ws, &msg); err != nil {
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
