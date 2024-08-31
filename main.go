package main

import (
	"fmt"
	"log"
	"net/http"

	"neon-chat/src"
)

func main() {
	log.Println("Application is starting...")
	src.SetupGlobalLogger(true, false)
	config := src.ReadEnvConfig()
	db := src.ConnectDB(config)
	app := src.InitAppState(config)

	src.SetupControllers(app, db)

	log.Printf("Listening at port [%d]\n", config.Port)
	runtineErr := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	log.Fatal(runtineErr)
}
