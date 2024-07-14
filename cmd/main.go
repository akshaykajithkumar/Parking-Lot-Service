package main

import (
	//"main/parking-lot-service/routes"
	"main/handlers"
	"main/models"
	"main/routes"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var db *gorm.DB

var server *handlers.Server

func initDB() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=youruser dbname=yourdb password=yourpassword")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.ParkingLot{}, &models.Tariff{}, &models.Ticket{})
	return db
}

func main() {
	db = initDB()
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	server = &handlers.Server{
		ParkingLots: make(map[string]*models.ParkingLot),
		//Clients:     make(map[*websocket.Conn]bool),
		Mutex: &sync.Mutex{},
	}
	routes.InitRoutes(e, server)

	e.Logger.Fatal(e.Start(":8080"))
}
