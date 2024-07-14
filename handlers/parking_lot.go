package handlers

import (
	"main/models"
	"net/http"
	"sync"

	//"github.com/gorilla/websocket"
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
	//Clients     map[*websocket.Conn]bool
	Mutex *sync.Mutex
}

func CreateParkingLot(c echo.Context) error {
	var parkingLot models.ParkingLot
	if err := c.Bind(&parkingLot); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if err := Db.Create(&parkingLot).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not create parking lot"})
	}

	for _, tariff := range parkingLot.Tariffs {
		tariff.ParkingLotID = parkingLot.ID
		if err := Db.Create(&tariff).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not save tariffs"})
		}
	}

	return c.JSON(http.StatusCreated, parkingLot)
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not retrieve parking lots"})
	}
	return c.JSON(http.StatusOK, parkingLots)
}
