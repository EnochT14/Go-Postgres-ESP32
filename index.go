package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux" // Import the Gorilla Mux router
	_ "github.com/lib/pq"
	"github.com/rs/cors" // Import the rs/cors package for CORS handling
)

const (
	host          = "89.168.96.111" // db.jisbrbnxxlltogqswsnx.supabase.co
	port          = 5432
	user          = "postgres"        //root
	password      = "u.!iQJ3itapJ9tp" //root , zq6G_gSJ:FeqErH
	dbname        = "postgres"        //esp32_data, test_db --- to change to postgres only db
	listenAddress = ":8080"
)

func main() {
	// Create a database connection
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ensure the database connection is alive
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create a Gorilla Mux router
	router := mux.NewRouter()

	// Add a CORS middleware with allowed origins
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Adjust as needed for your requirements
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Attach the CORS middleware to the router
	router.Use(corsMiddleware.Handler)

	// Create an HTTP endpoint to receive sensor data
	router.HandleFunc("/collect-sensor-data", func(w http.ResponseWriter, r *http.Request) {
		// Parse the JSON data from the request body
		var sensorData struct {
			Temperature float64 `json:"temperature"`
			Humidity    float64 `json:"humidity"`
			Pressure    float64 `json:"pressure"`
			CO2PPM      int     `json:"co2_ppm"`
			TVOCPpb     int     `json:"tvoc_ppb"`
		}

		if err := json.NewDecoder(r.Body).Decode(&sensorData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get the current timestamp
		timestamp := time.Now()

		// Insert data into the sensor_data table
		insertStatement := `
            INSERT INTO sensor_data (timestamp, temperature, humidity, pressure, co2_ppm, tvoc_ppb)
            VALUES ($1, $2, $3, $4, $5, $6)
        `
		_, err := db.Exec(insertStatement, timestamp, sensorData.Temperature, sensorData.Humidity, sensorData.Pressure, sensorData.CO2PPM, sensorData.TVOCPpb)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Sensor data inserted successfully.")
		w.WriteHeader(http.StatusNoContent)
	})

	// Start the HTTP server with the router
	log.Printf("Server listening on %s...", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, router))
}
