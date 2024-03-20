package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.chat/controller"
	"go.chat/db"
	"go.chat/model"
	"go.chat/utils"
)

func main() {
	// log setup
	now := time.Now()
	timestamp := now.Format(time.RFC3339)
	date := strings.Split(timestamp, "T")[0]
	logPath := fmt.Sprintf("log/from-%s.log", date)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	// write log to both file and stderr
	multi := io.MultiWriter(logFile, os.Stderr)
	log.SetOutput(multi)
	// parse args
	args, err := utils.ArgsRead()
	if err != nil {
		log.Printf("Error parsing args: %v\n", err)
		log.Println(utils.ArgsHelp())
		os.Exit(13)
	}
	log.Printf("  args: %v\n", *args)
	// TODO args.DBPath
	dbPath := "db/chat.db"
	db, err := db.ConnectDB(dbPath)
	if err != nil {
		log.Fatalf("Error opening db: %s", err)
	}
	log.Println("Setting up application")
	app := &model.ApplicationState
	log.Println("Setting up controllers")
	controller.Setup(app, db)
	log.Printf("Starting server at port [%d]\n", args.Port)
	runtineErr := http.ListenAndServe(fmt.Sprintf(":%d", args.Port), nil)
	log.Fatal(runtineErr)
}
