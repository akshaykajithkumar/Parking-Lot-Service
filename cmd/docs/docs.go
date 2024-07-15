// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/parking-lots": {
            "get": {
                "description": "Get a list of all parking lots",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "parkinglots"
                ],
                "summary": "List all parking lots",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.ParkingLotSummary"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new parking lot [first_rate is the rate of first hours(eg : 1st day = 24 hours ) and after_rate is the rate of after hours]",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "parkinglots"
                ],
                "summary": "Create Parking Lot",
                "parameters": [
                    {
                        "description": "Parking Lot",
                        "name": "parkinglot",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SwaggerParkingLot"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.SwaggerParkingLot"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/parking-lots/{id}": {
            "get": {
                "description": "Get detailed information about a specific parking lot",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "parkinglots"
                ],
                "summary": "Get details of a parking lot",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Parking Lot ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SwaggerParkingLot"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/parking-lots/{id}/available-spots": {
            "get": {
                "description": "Get available spots in a specific parking lot",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "parkinglots"
                ],
                "summary": "Get available spots",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Parking Lot ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AvailableSpotsResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/parking-lots/{id}/park": {
            "post": {
                "description": "Park a vehicle in a specific parking lot",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "parkinglots"
                ],
                "summary": "Park a vehicle",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Parking Lot ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Vehicle Data",
                        "name": "vehicle",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.VehicleData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SwaggerTicket"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/parking-lots/{id}/unpark": {
            "post": {
                "description": "Unpark a vehicle from a specific parking lot",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "parkinglots"
                ],
                "summary": "Unpark a vehicle",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Parking Lot ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Vehicle Data",
                        "name": "vehicle",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UnparkVehicleData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UnparkResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AvailableSpotsResponse": {
            "type": "object",
            "properties": {
                "bus": {
                    "type": "integer"
                },
                "car": {
                    "type": "integer"
                },
                "motorcycle": {
                    "type": "integer"
                }
            }
        },
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "models.ParkingLotSummary": {
            "type": "object",
            "properties": {
                "ID": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models.Receipt": {
            "type": "object",
            "properties": {
                "ID": {
                    "type": "integer"
                },
                "entry_time": {
                    "type": "string"
                },
                "exit_time": {
                    "type": "string"
                },
                "parking_lot_id": {
                    "type": "integer"
                },
                "vehicle_number": {
                    "type": "string"
                },
                "vehicle_type": {
                    "type": "string"
                }
            }
        },
        "models.SwaggerParkingLot": {
            "type": "object",
            "properties": {
                "bus_spots": {
                    "type": "integer",
                    "example": 20
                },
                "car_spots": {
                    "type": "integer",
                    "example": 200
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "motorcycle_spots": {
                    "type": "integer",
                    "example": 50
                },
                "name": {
                    "type": "string",
                    "example": "Parking lot A"
                },
                "occupied_buses": {
                    "type": "integer",
                    "example": 5
                },
                "occupied_cars": {
                    "type": "integer",
                    "example": 150
                },
                "occupied_motorcycles": {
                    "type": "integer",
                    "example": 10
                },
                "tariffs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.SwaggerTariff"
                    }
                }
            }
        },
        "models.SwaggerRatePlan": {
            "type": "object",
            "properties": {
                "after_rate": {
                    "type": "number",
                    "example": 5
                },
                "first_hours": {
                    "type": "integer",
                    "example": 2
                },
                "first_rate": {
                    "type": "number",
                    "example": 10
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "tariff_id": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "models.SwaggerTariff": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "parking_lot_id": {
                    "type": "integer",
                    "example": 1
                },
                "rate_plans": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.SwaggerRatePlan"
                    }
                },
                "vehicle_type": {
                    "type": "string",
                    "example": "car"
                }
            }
        },
        "models.SwaggerTicket": {
            "type": "object",
            "properties": {
                "ID": {
                    "type": "integer"
                },
                "entry_time": {
                    "type": "string"
                },
                "parking_lot_name": {
                    "type": "string"
                },
                "vehicle_number": {
                    "type": "string"
                },
                "vehicle_type": {
                    "type": "string"
                }
            }
        },
        "models.UnparkResponse": {
            "type": "object",
            "properties": {
                "Total Fee": {
                    "type": "integer"
                },
                "receipt": {
                    "$ref": "#/definitions/models.Receipt"
                }
            }
        },
        "models.UnparkVehicleData": {
            "type": "object",
            "properties": {
                "vehicle_number": {
                    "type": "string"
                }
            }
        },
        "models.VehicleData": {
            "type": "object",
            "properties": {
                "vehicle_number": {
                    "type": "string"
                },
                "vehicle_type": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8181",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Parking Lot API",
	Description:      "This is a parking lot management API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
