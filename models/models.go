package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ParkingLot struct {
	gorm.Model
	Name                string   `json:"name"`
	MotorcycleSpots     int      `json:"motorcycle_spots"`
	CarSpots            int      `json:"car_spots"`
	BusSpots            int      `json:"bus_spots"`
	OccupiedMotorcycles int      `json:"occupied_motorcycles"`
	OccupiedCars        int      `json:"occupied_cars"`
	OccupiedBuses       int      `json:"occupied_buses"`
	Tariffs             []Tariff `json:"tariffs" gorm:"foreignKey:ParkingLotID"`
}

type Tariff struct {
	gorm.Model
	ParkingLotID uint       `json:"parking_lot_id"`
	VehicleType  string     `json:"vehicle_type"`
	RatePlans    []RatePlan `json:"rate_plans" gorm:"foreignKey:TariffID"`
}

type RatePlan struct {
	gorm.Model
	TariffID   uint    `json:"tariff_id"`
	FirstHours int     `json:"first_hours"`
	FirstRate  float64 `json:"first_rate"`
	AfterRate  float64 `json:"after_rate"`
}

type Ticket struct {
	gorm.Model
	VehicleType   string     `json:"vehicle_type"`
	VehicleNumber string     `json:"vehicle_number"`
	ParkingLotID  uint       `json:"parking_lot_id"`
	SpotNumber    uint       `json:"spot_number"`
	EntryTime     time.Time  `json:"entry_time"`
	ExitTime      *time.Time `json:"exit_time,omitempty"`
	ParkingLot    ParkingLot `gorm:"foreignkey:ParkingLotID" json:"-"`
}
