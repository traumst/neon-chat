package main

import (
	"fmt"
	"log"
	"net/http"

	"neon-chat/src"
	test "neon-chat/src/_test"
	"neon-chat/src/db"
	"neon-chat/src/utils/config"
)

func main() {
	log.Println("Starting up...")
	src.SetupGlobalLogger(true, false)
	config := src.ReadEnvConfig()
	log.Println("Connecting db...")
	db := src.ConnectDB(config.Sqlite)
	log.Println("Verifying db requirements...")
	initTestData(db, config.TestDataInsert, config.TestUsers)
	log.Println("Creating state...")
	app := src.InitAppState(config)
	log.Println("Setting up controllers...")
	src.SetupControllers(app, db)
	log.Printf("Application is starting at port [%d]\n", config.Port)
	runtineErr := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	log.Fatal(runtineErr)
}

func initTestData(db *db.DBConn, insertTestData bool, testUsers config.TestUsers) {
	if !insertTestData {
		log.Println("Skipping test data insert")
		return
	}
	log.Println("Inserting test users...")
	userCount, err := test.CreateTestUsers(db, testUsers)
	if err != nil || userCount < 0 {
		log.Fatalf("ERROR failed to create [%d] of [%d] test users: %s", userCount, len(testUsers), err)
	} else if userCount == 0 {
		log.Printf("created none of [%d] specified test users", len(testUsers))
	} else /* if userCount > 0 */ {
		log.Printf("created [%d] out of [%d] specified test users", userCount, len(testUsers))
	}
	log.Println("Inserting test auth...")
	authCount, err := test.CreateTestAuth(db, testUsers)
	if err != nil || authCount < 0 {
		log.Fatalf("ERROR failed to create test auth: %s", err)
	} else if authCount == 0 {
		log.Println("created none of specified test auth")
	} else /* if authCount > 0 */ {
		log.Println("created specified test auth")
	}
	if userCount != authCount {
		log.Fatalf("ERROR created users count [%d] does not match created auth count [%d]", userCount, authCount)
	}
}
