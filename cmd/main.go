package main

import (
	"log"
	"main/cmd/docs"
	"main/handlers"
	"main/models"
	"main/routes"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func main() {
	docs.SwaggerInfo.Title = "Parking lot service"
	docs.SwaggerInfo.Host = "localhost:8181"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http"}
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=parkingservice password=admin sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	db.AutoMigrate(&models.ParkingLot{}, &models.Tariff{}, &models.Ticket{}, &models.RatePlan{})

	handlers.SetDB(db)

	server := &handlers.Server{
		ParkingLots: make(map[string]*models.ParkingLot),
		Clients:     make(map[*websocket.Conn]bool),
		Mutex:       &sync.Mutex{},
		DB:          db,
	}

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.InitRoutes(e, server)

	e.Logger.Fatal(e.Start(":8181"))
}
