package routes

import (
	"main/handlers"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Parking Lot API
// @version 1.0
// @description This is a parking lot management API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8181
// @BasePath /api/v1
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
	// e.POST("/api/v1/parkinglot", handlers.CreateParkingLot)
	// e.GET("/api/v1/parkinglot/:id", handlers.GetParkingLotDetails)
	// e.GET("/api/v1/parkinglots", handlers.ListParkingLots)
	// e.POST("/api/v1/parkinglot/:id/park", handlers.ParkVehicle)
	// e.POST("/api/v1/parkinglot/:id/unpark", handlers.UnparkVehicle)
	// e.GET("/api/v1/parkinglot/:id/availablespots", handlers.GetAvailableSpots)

	// Swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
