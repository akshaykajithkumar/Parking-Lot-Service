package main

import (
	"flag"
	"main/client"
	"main/handlers"
	"main/models"
	"main/routes"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/net/websocket"
)

var Db *gorm.DB
var server *handlers.Server

func initDB() *gorm.DB {
	// Update the connection string as per your database configuration
	Db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=parkingservice password=admin sslmode=disable")

	if err != nil {
		panic("failed to connect database")
	}
	Db.AutoMigrate(&models.ParkingLot{}, &models.Tariff{}, &models.Ticket{})
	return Db
}

func main() {
	runClient := flag.Bool("client", false, "Run the WebSocket client")
	flag.Parse()

	if *runClient {
		client.RunWebSocketClient()
		return
	}

	Db = initDB()
	defer Db.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	server = &handlers.Server{
		ParkingLots: make(map[string]*models.ParkingLot),
		Clients:     make(map[*websocket.Conn]bool),
		Mutex:       &sync.Mutex{},
	}

	// Serve static files
	e.Static("/", "static")

	routes.InitRoutes(e, server)

	e.Logger.Fatal(e.Start(":8181"))
}
