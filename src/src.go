package src

import (
	"fmt"
	"io"
	"log"
	"os"
	"prplchat/src/db"
	"prplchat/src/handler/state"
	"prplchat/src/utils"
	"strings"
	"time"
)

func SetupGlobalLogger() {
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
}

func ReadEnvConfig() *utils.Config {
	log.Println("parsing config...")
	//config, err := utils.ArgsRead()
	config, err := utils.EnvRead()
	if err != nil {
		log.Printf("Error parsing config: %s\n", err)
		log.Println(utils.ConfigHelp())
		os.Exit(13)
	}
	log.Printf("\tparsed config: %s\n", config)
	return config
}

func ConnectDB(config *utils.Config) *db.DBConn {
	log.Println("connecting db...")
	db, err := db.ConnectDB(config.Sqlite)
	if err != nil {
		log.Fatalf("Error opening db at [%s]: %s", config.Sqlite, err)
	}
	return db
}

func InitAppState(config *utils.Config) *state.State {
	log.Println("init app state...")
	app := &state.GlobalAppState
	app.Init(utils.Config{
		CacheSize: config.CacheSize,
		Port:      config.Port,
		Sqlite:    config.Sqlite,
		Smtp: utils.SmtpConfig{
			User: config.Smtp.User,
			Pass: config.Smtp.Pass,
			Host: config.Smtp.Host,
			Port: config.Smtp.Port,
		},
	})
	return app
}
