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

func ReadEnvConfig() *config.Config {
	log.Println("parsing config...")
	c, err := config.EnvRead()
	if err != nil {
		log.Printf("Error parsing config: %s\n", err)
		log.Println(config.ConfigHelp())
		os.Exit(13)
	}
	log.Printf("\tparsed config: %s\n", c)
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

func CreateTestUsers(dbConn *db.DBConn, TestUsers config.TestUsers) (int, error) {
	log.Println("Checking test users status...")
	dbUsers, err := db.SearchUsers(dbConn.Conn, TestUsers.GetNames()...)
	if err != nil {
		return -1, fmt.Errorf("failed to search for test users: %s", err)
	}
	newUsers := make(config.TestUsers, 0)
	for _, testUser := range TestUsers {
		exists := false
		for _, dbUser := range dbUsers {
			if dbUser.Name == testUser.Name {
				exists = true
				break
			}
		}
		if !exists {
			newUsers = append(newUsers, testUser)
		}
	}
	log.Printf("there are [%d] test users to create: \n%+v\n", len(newUsers), newUsers)
	// for _, testUser := range newUsers {
	// 	salt := utils.GenerateSalt(testUser.Name, string(app.LocalUserType))
	// 	dbUser := db.User{
	// 		Id:     0,
	// 		Name:   testUser.Name,
	// 		Email:  testUser.Email,
	// 		Type:   string(enum.AuthTypeEmail),
	// 		Status: "",
	// 		Salt:   salt,
	// 	}

	// 	newDbUser, err := db.AddUser(dbConn.Tx, &dbUser)
	// 	if err != nil {
	// 		return -999, fmt.Errorf("failed to create test user[%s]: %s", testUser.Name, err)
	// 	}
	// }

	return -999, nil
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
