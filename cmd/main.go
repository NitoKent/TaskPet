package main

import (
	"log"
	"os"
	"strconv"
	"taskPet/m/v2/cmd/api"
	"taskPet/m/v2/db"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	cfg := db.Config{
		Username: dbUser,
		Password: dbPassword,
		Host:     dbHost,
		Port:     dbPort,
		SSLMode:  "disable",
	}

	database, err := db.NewStorage(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer database.Close()

	server := api.NewAPIServer(":8056", database)
	if err := server.Run(); err != nil {
		log.Fatalf("Server start failed: %v", err)
	}
}
