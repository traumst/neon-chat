package src

import (
	"fmt"
	"io"
	"log"
	"neon-chat/src/db"
	"neon-chat/src/state"
	"neon-chat/src/utils/config"
	"os"
	"strings"
	"time"
)

func SetupGlobalLogger(stderr bool, dir string) {
	log.Printf("TRACE SetupGlobalLogger stderr[%t], dir[%s]", stderr, dir)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)
	var w io.Writer
	var logFile *os.File
	if len(dir) > 0 && stderr {
		log.Println("INFO Logging to file and stderr")
		logFile = OpenLogFile(dir)
		w = io.MultiWriter(os.Stderr, logFile)
	} else if len(dir) > 0 {
		log.Println("INFO Logging to file")
		logFile = OpenLogFile(dir)
		w = logFile
	} else if stderr {
		log.Println("INFO Logging to stderr")
		w = os.Stderr
	} else {
		log.Println("WARN Logging is DISABLED")
		w = io.Discard
	}
	if logFile != nil {
		defer logFile.Close()
	}
	log.SetPrefix("\n\t")
	log.SetOutput(w)
}

func OpenLogFile(dir string) *os.File {
	timestamp := time.Now().Format(time.RFC3339)
	date := strings.Split(timestamp, "T")[0]
	logPath := fmt.Sprintf("%s/from-%s.log", dir, date)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return logFile
}

func ReadEnvConfig() *config.Config {
	log.Println("parsing config...")
	c, err := config.EnvRead()
	if err != nil {
		log.Printf("Error parsing config: %s\n", err)
		log.Println(config.ConfigHelp())
		os.Exit(13)
	}
	//log.Printf("...parsed config\n%s\n", c)
	return c
}

func ConnectDB(dbFilePath string) *db.DBConn {
	log.Println("connecting db...")
	db, err := db.ConnectDB(dbFilePath)
	if err != nil {
		log.Fatalf("Error opening db at [%s]: %s", dbFilePath, err)
	}
	return db
}

func InitAppState(c *config.Config) *state.State {
	log.Println("init app state...")
	app := &state.GlobalAppState
	app.Init(config.Config{
		CacheSize: c.CacheSize,
		Port:      c.Port,
		Sqlite:    c.Sqlite,
		Smtp: config.SmtpConfig{
			User: c.Smtp.User,
			Pass: c.Smtp.Pass,
			Host: c.Smtp.Host,
			Port: c.Smtp.Port,
		},
		TestUsers: c.TestUsers[:],
	})
	return app
}
