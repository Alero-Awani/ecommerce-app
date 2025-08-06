package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort            string
	Dsn                   string
	AppSecret             string
	TwilioAccountSid      string
	TwilioAuthToken       string
	TwilioFromPhoneNumber string
}

// function to read environment variables and return application struct

func SetupEnv() (cfg AppConfig, err error) {
	log.Println("Getting Environment Variables from configmap")

	if os.Getenv("APP_ENV") == "dev" {
		err = godotenv.Load(".env")
		log.Fatalf("Error loading .env file: %v\n", err)
	}
	// log.Println("Loading http port from environment variables\n")
	// httpPort := os.Getenv("HTTP_PORT")

	// if len(httpPort) < 1 {
	// 	return AppConfig{}, errors.New("env variables not found")
	// }

	// Dev
	// Dsn := os.Getenv("DSN")
	// if len(Dsn) < 1 {
	// 	return AppConfig{}, errors.New("env variables not found")
	// }

	// Production
	log.Print("Loading database connection string from environment variables configmaps\n")
	
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	log.Println("The Database host is:", dbHost)

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		return AppConfig{}, fmt.Errorf("required database environment variables not found (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)")
	}

	fmt.Printf("Database connection string: host=%s user=%s dbname=%s port=%s\n", 
		dbHost, dbUser, dbName, dbPort) // Don't log password

	Dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// appSecret := os.Getenv("APP_SECRET")
	// if len(appSecret) < 1 {
	// 	return AppConfig{}, errors.New("env variables not found")
	// }

	return AppConfig{
		// ServerPort: httpPort, 
		Dsn: Dsn, 
		// AppSecret: appSecret,
		TwilioAccountSid:      os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken:       os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioFromPhoneNumber: os.Getenv("TWILIO_FROM_PHONE_NUMBER"),
	}, nil
}
