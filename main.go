package main

import (
	"auth/router"
	"fmt"
	"github.com/rs/cors"
	"log"
)

func main() {

	t := router.Router()

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	corsWrapper.Handler(t)

	runError := t.Run(":8080")
	if runError != nil {
		log.Fatalln("failed to run app")
	}
	fmt.Println("Starting server on the port 8080...")
}
