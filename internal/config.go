package internal

import (
	"fmt"
	"log"
	"os"
)

type MainConfig struct {
	AnalysisServiceURL string
	DbConnectionString string
}

func MakeMainConfigFromENV() *MainConfig {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	analysisServiceURL := os.Getenv("ANALYSIS_SERVICE_URL")

	// Validate environment variables
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || analysisServiceURL == "" {
		log.Fatal("Missing required environment variables")
	}

	// Construct the database connection string
	dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	return &MainConfig{
		analysisServiceURL,
		dbConnectionString,
	}
}
