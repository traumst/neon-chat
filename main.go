package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"neon-chat/src"
	test "neon-chat/src/_test"
	"neon-chat/src/db"
	"neon-chat/src/state"
	"neon-chat/src/utils"
	"neon-chat/src/utils/config"
	h "neon-chat/src/utils/http"
)

func main() {
	log.Println("Reading config...")
	config := src.ReadEnvConfig()
	log.Println("Setup logger...")
	src.SetupGlobalLogger(config.Log.Stdout, config.Log.Dir)
	log.Println("Connecting db...")
	db := src.ConnectDB(config.Sqlite)
	log.Println("Verifying db requirements...")
	initTestData(db, config.TestDataInsert, config.TestUsers)
	// TODO ensure maintenance completes
	go db.ScheduleMaintenance()
	//
	log.Println("Creating state...")
	app := src.InitAppState(config)
	//
	log.Println("Setting up controllers...")
	src.SetupControllers(app, db, config.RateLimit)
	//
	log.Println("Loading previous sessions...")
	if err := h.LoadSessionsFromFile(config.BackupConfig.SessionFilePath); err != nil {
		log.Printf("Could not load sessions: %v", err)
	}

	log.Printf("Application is starting at port [%d]\n", config.Port)
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Port),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		log.Fatalf("Server stopping [%s], %s\n", server.Addr, err)
	}()

	<-stop

	log.Println("Shutting down gracefully, press Ctrl+C again to force")
	gracefulShutdown(app, db, server, config)

	log.Println("Server stopped")
}

func gracefulShutdown(
	app *state.State,
	db *db.DBConn,
	server *http.Server,
	conf *config.Config,
) {
	log.Println("Waiting for users to leave...")
	err := utils.MaintenanceManager.RaiseFlag()
	if err != nil {
		log.Fatalf("ERROR failed to raise maintenance flag: %s", err)
	}
	activeCount := utils.MaintenanceManager.WaitUsersLeave(3 * time.Second)
	if activeCount != 0 {
		log.Printf("ERROR [%d] users did not leave", activeCount)
	}
	//
	log.Println("Storing existing sessions...")
	if err := h.SaveSessionsToFile(conf.BackupConfig.SessionFilePath); err != nil {
		log.Printf("Could not save sessions: %v", err)
	}
	//
	log.Println("Storing open user chats...")
	if err := app.SaveToFile(conf.BackupConfig.UserChatFilePath); err != nil {
		log.Printf("Could not save user chat: %v", err)
	}
	//
	log.Println("Closing db connection...")
	db.ConnClose(3 * time.Second)
	//
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
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
