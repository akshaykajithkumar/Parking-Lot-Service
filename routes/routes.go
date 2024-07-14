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

}
