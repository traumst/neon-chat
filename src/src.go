package src

import (
	"fmt"
	"io"
	"log"
	"neon-chat/src/db"
	"neon-chat/src/handler/state"
	"neon-chat/src/utils"
	"os"
	"strings"
	"time"
)

func SetupGlobalLogger(toStderr bool, toFile bool) {
	log.Println("TRACE IN SetupGlobalLogger")
	now := time.Now()
	timestamp := now.Format(time.RFC3339)
	date := strings.Split(timestamp, "T")[0]
	logPath := fmt.Sprintf("log/from-%s.log", date)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// log.SetPrefix("CUSTOM_LOG: ")
	if toFile && toStderr {
		log.SetOutput(io.MultiWriter(logFile, os.Stderr))
		log.Printf("Logging to file [%s] and stderr\n", logPath)
	} else if toFile {
		log.SetOutput(logFile)
		log.Printf("Logging to file [%s]\n", logPath)
	} else {
		log.SetOutput(os.Stderr)
		log.Println("Logging to stderr")
	}
	log.Println("TRACE OUT SetupGlobalLogger")
}

func ReadEnvConfig() *utils.Config {
	log.Println("parsing config...")
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
