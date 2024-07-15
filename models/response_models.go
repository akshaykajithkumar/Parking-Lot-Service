package models

import (
	"time"
)

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}
type VehicleData struct {
	VehicleType   string `json:"vehicle_type"`
	VehicleNumber string `json:"vehicle_number"`
}

// UnparkVehicleData represents the request body for unparking a vehicle
type UnparkVehicleData struct {
	VehicleNumber string `json:"vehicle_number"`
}

// AvailableSpotsResponse represents the response for available spots
type AvailableSpotsResponse struct {
	Motorcycle int `json:"motorcycle"`
	Car        int `json:"car"`
	Bus        int `json:"bus"`
}

// SwaggerParkingLot is used to generate Swagger documentation
type SwaggerParkingLot struct {
	ID uint `json:"id" example:"1"`

	Name                string          `json:"name" example:"Parking lot A"`
	MotorcycleSpots     int             `json:"motorcycle_spots" example:"50"`
	CarSpots            int             `json:"car_spots" example:"200"`
	BusSpots            int             `json:"bus_spots" example:"20"`
	OccupiedMotorcycles int             `json:"occupied_motorcycles" example:"10"`
	OccupiedCars        int             `json:"occupied_cars" example:"150"`
	OccupiedBuses       int             `json:"occupied_buses" example:"5"`
	Tariffs             []SwaggerTariff `json:"tariffs"`
}

// SwaggerTariff is used to generate Swagger documentation
type SwaggerTariff struct {
	ID           uint              `json:"id" example:"1"`
	ParkingLotID uint              `json:"parking_lot_id" example:"1"`
	VehicleType  string            `json:"vehicle_type" example:"car"`
	RatePlans    []SwaggerRatePlan `json:"rate_plans"`
}

// SwaggerRatePlan is used to generate Swagger documentation
type SwaggerRatePlan struct {
	ID         uint    `json:"id" example:"1"`
	TariffID   uint    `json:"tariff_id" example:"1"`
	FirstHours int     `json:"first_hours" example:"2"`
	FirstRate  float64 `json:"first_rate" example:"10.00"`
	AfterRate  float64 `json:"after_rate" example:"5.00"`
}

// SwaggerTicket represents the parking ticket response structure
type SwaggerTicket struct {
	ID             int       `json:"ID"`
	EntryTime      time.Time `json:"entry_time"`
	ParkingLotName string    `json:"parking_lot_name"`
	VehicleNumber  string    `json:"vehicle_number"`
	VehicleType    string    `json:"vehicle_type"`
}

// Receipt represents the receipt structure for the unparked vehicle
type Receipt struct {
	ID            int       `json:"ID"`
	EntryTime     time.Time `json:"entry_time"`
	ExitTime      time.Time `json:"exit_time"`
	ParkingLotID  int       `json:"parking_lot_id"`
	VehicleNumber string    `json:"vehicle_number"`
	VehicleType   string    `json:"vehicle_type"`
}

// UnparkResponse represents the response structure for the unpark action
type UnparkResponse struct {
	TotalFee int     `json:"Total Fee"`
	Receipt  Receipt `json:"receipt"`
}

// ParkingLotSummary represents a summary of a parking lot
type ParkingLotSummary struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
}
