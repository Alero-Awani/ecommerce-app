package main

import (
	"ecommerce-app/config"
	"ecommerce-app/internal/api"
	"log"
)

func main() {
	cfg, err := config.SetupEnv()
	log.Printf("Config loaded: %+v\n", cfg)

	if err != nil {
		log.Fatalf("This config file is not loaded properly %v\n", err)
		log.Printf("The error message is:", err)
	}

	api.StartServer(cfg)
}
