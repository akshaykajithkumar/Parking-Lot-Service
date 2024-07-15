package main

import (
	"log"
	"main/handlers"
	"main/models"
	"main/routes"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.InitRoutes(e, server)

	e.Logger.Fatal(e.Start(":8181"))
}
