package routes

import (
	"main/handlers"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, server *handlers.Server) {
	handlers.SetDB(handlers.Db)
	handlers.InitServer(server)

	e.POST("/parking-lots", handlers.CreateParkingLot)
	e.GET("/parking-lots/:id", handlers.GetParkingLotDetails)
	e.GET("/parking-lots", handlers.ListParkingLots)
	e.POST("/parking-lots/:id/park", handlers.ParkVehicle)
	e.POST("/parking-lots/:id/unpark", handlers.UnparkVehicle)
	e.GET("/parking-lots/:id/available-spots", handlers.GetAvailableSpots)
	e.GET("/ws", handlers.HandleWebSocket)
}
