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
	"neon-chat/src/db"
	"neon-chat/src/state"
	"neon-chat/src/utils"
	"neon-chat/src/utils/config"
	h "neon-chat/src/utils/http"
)

func main() {
	app, db, server, conf := serverStartup()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		log.Fatalf("Server stopping [%s], %s\n", server.Addr, err)
	}()

	<-stop

	log.Println("Shutting down gracefully, press Ctrl+C again to force")
	gracefulShutdown(app, db, server, conf)

	log.Println("Server stopped")
}

func serverStartup() (
	app *state.State,
	db *db.DBConn,
	server *http.Server,
	conf *config.Config,
) {
	log.Println("Reading config...")
	conf = src.ReadEnvConfig()
	log.Println("Setup logger...")
	src.SetupGlobalLogger(conf.Log.Stdout, conf.Log.Dir)
	log.Println("Initiating db...")
	db = src.InitDBConn(conf.Sqlite, conf.TestDataInsert, conf.TestUsers)
	log.Println("Initiating app state...")
	app = src.InitAppState(conf)
	log.Println("Setting up controllers...")
	src.SetupControllers(app, db, conf.RateLimit)
	log.Println("Loading previous sessions...")
	if err := h.LoadSessionsFromFile(conf.BackupConfig.SessionFilePath); err != nil {
		log.Printf("Could not load sessions: %s", err)
	}
	server = &http.Server{
		Addr: fmt.Sprintf(":%d", conf.Port),
	}
	return app, db, server, conf
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
		log.Printf("ERROR failed to raise maintenance flag: %s", err)
		res := utils.MaintenanceManager.WaitMaintenanceComplete(30 * time.Second)
		if !res {
			log.Printf("ERROR server stuck in maintenance, data may become corrupted")
		} else {
			log.Printf("INFO no maintenance now, safe to shutdown")
		}
		_ = utils.MaintenanceManager.RaiseFlag()
	}
	activeCount := utils.MaintenanceManager.WaitUsersLeave(3 * time.Second)
	if activeCount != 0 {
		log.Printf("ERROR [%d] users did not leave", activeCount)
	}
	log.Println("Storing existing sessions...")
	if err := h.SaveSessionsToFile(conf.BackupConfig.SessionFilePath); err != nil {
		log.Printf("Could not save sessions: %s", err)
	}
	log.Println("Storing open user chats...")
	if err := app.SaveToFile(conf.BackupConfig.UserChatFilePath); err != nil {
		log.Printf("Failed to save user chats: %s", err)
	}
	log.Println("Closing db connection...")
	if err := db.ConnClose(3 * time.Second); err != nil {
		log.Printf("Failed to close db connection: %s", err)
	}
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}
}
