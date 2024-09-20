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
	db := src.ConnectDB(config.Sqlite)
	log.Println("Verifying db requirements...")
	res, err := src.CreateTestUsers(db, config.TestUsers)
	if err != nil {
		log.Fatalf("ERROR creating specified test users: %s", err)
	} else if res < 0 {
		log.Fatalf("ERROR creating specified test users, err_code[%d]", res)
	} else if res == 0 {
		log.Printf("created none of [%d] specified test users", len(config.TestUsers))
	} else /* if res > 0 */ {
		log.Printf("created [%d] out of [%d] specified test users", res, len(config.TestUsers))
	}
	log.Println("Creating state...")
	app := src.InitAppState(config)
	log.Println("Setting up controllers...")
	src.SetupControllers(app, db)
	log.Printf("Application is starting at port [%d]\n", config.Port)
	runtineErr := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	log.Fatal(runtineErr)
}
