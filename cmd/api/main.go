package main

import (
	"log"

	"github.com/charmingruby/backpago/configs"
	"github.com/charmingruby/backpago/pkg/database"
)

func main() {
	configs.LoadConfigs()

	_, err := database.NewConnection()
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
}
