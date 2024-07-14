package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ParkingLot struct {
	gorm.Model
	Name                string
	MotorcycleSpots     int
	CarSpots            int
	BusSpots            int
	OccupiedMotorcycles int
	OccupiedCars        int
	OccupiedBuses       int
	Tariffs             []Tariff
}

type Tariff struct {
	gorm.Model
	ParkingLotID uint
	VehicleType  string
	RatePerHour  float64
	MaxDailyRate float64
}

type Ticket struct {
	gorm.Model
	VehicleType  string
	LicensePlate string
	ParkingLotID uint
	SpotNumber   uint
	EntryTime    time.Time
}
