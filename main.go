package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"prplchat/src/controller"
	"prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/utils"
)

func main() {
	log.Println("Application is starting...")
	// init log
	log.Println("setting up logger...")
	now := time.Now()
	timestamp := now.Format(time.RFC3339)
	date := strings.Split(timestamp, "T")[0]
	logPath := fmt.Sprintf("log/from-%s.log", date)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	multi := io.MultiWriter(logFile, os.Stderr)
	log.SetOutput(multi)

	// parse args
	log.Println("parsing config...")
	//config, err := utils.ArgsRead()
	config, err := utils.EnvRead()
	if err != nil {
		log.Printf("Error parsing config: %s\n", err)
		log.Println(utils.Help())
		os.Exit(13)
	}
	log.Printf("\tparsed config: %s\n", config)

	log.Println("connecting db...")
	db, err := db.ConnectDB(config.Sqlite)
	if err != nil {
		log.Fatalf("Error opening db at [%s]: %s", config.Sqlite, err)
	}

	log.Println("init app state...")
	app := &handler.ApplicationState
	app.Init(db, utils.Config{
		LoadLocal: config.LoadLocal,
		Smtp: utils.SmtpConfig{
			User: config.Smtp.User,
			Pass: config.Smtp.Pass,
			Host: config.Smtp.Host,
			Port: config.Smtp.Port,
		},
	})

	log.Println("init controllers...")
	controller.Setup(app, db)

	log.Printf("Starting server at port [%d]\n", config.Port)
	runtineErr := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	log.Fatal(runtineErr)
}
