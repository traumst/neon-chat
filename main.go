package main

import (
	"fmt"
	"log"
	"net/http"

	"neon-chat/src"
)

func main() {
	log.Println("Starting up...")
	src.SetupGlobalLogger(true, false)
	config := src.ReadEnvConfig()
	log.Println("Connecting db...")
	db := src.ConnectDB(config)
	log.Println("Creating state...")
	app := src.InitAppState(config)
	log.Println("Setting up controllers...")
	src.SetupControllers(app, db)
	log.Printf("Application is starting at port [%d]\n", config.Port)
	runtineErr := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	log.Fatal(runtineErr)
}
