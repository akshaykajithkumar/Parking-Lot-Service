

# Parking Lot Service
## Features

- **Parking Lot Management**: Create, list, and retrieve details of parking lots.
- **Real-time Updates**: Monitor real-time availability of parking spots.
- **Dynamic Parking Fee**: Implement dynamic pricing for parking based on time .
- **Vehicle Management**: Park and unpark vehicles with ease.


## Getting Started

### Prerequisites

- Docker (for containerization)
- Postman (for testing WebSocket connection)

### Running the Application

1. **Clone the repository**:

```sh
git clone https://github.com/akshaykajithkumar/Parking-Lot-Service

cd Parking Lot Service

```

### Run the Application
 
```sh
go run cmd/main.go 
```


2. **Build and run the Docker container**:
```sh
docker build -t parking-lot-service .

```

```sh
docker run -p 8181:8181 parking-lot-service

```
3. **Access the Swagger UI**:
    Open your browser and navigate to [http://localhost:8181/swagger/index.html](http://localhost:8181/swagger/index.html) to explore and test the API endpoints.

### Testing the WebSocket Connection

1. Open Postman and create a new WebSocket request.
 

2. Connect to the WebSocket endpoint:
    ```
    ws://localhost:8181/ws
    ```
3. To get the number of available parking spots for a specific parking lot, use the following HTTP endpoint: then visit the WebSocket endpoint for monitoring the data
    ```
    http://localhost:8181/parking-lots/{id}/available-spots
    ```
    Replace `{id}` with the actual parking lot ID you want to query.

## Technologies Used

- **Echo Framework**: Used for building the HTTP server and handling HTTP requests.
- **PostgreSQL**: Database system used for storing parking lot and related data.
- **WebSocket**: Utilized for real-time communication to provide instant updates on parking spot availability.
- **GORM**: Object-Relational Mapping (ORM) library for interfacing with PostgreSQL, simplifying database operations.
