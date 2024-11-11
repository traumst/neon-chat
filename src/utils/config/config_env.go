package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func ConfigHelp() string {
	return `Application expects the config file '.env' in the root directory.
		* find .env.template
		* copy it to .env
		* set desired values`
}

func EnvRead() (*Config, error) {
	envFileRootPath := ".env"
	envFile, err := os.OpenFile(envFileRootPath, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open .env file from [%s]: %v", envFileRootPath, err)
	}

	buffer := make([]byte, 1024)
	n, err := envFile.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read .env file from [%s]: %v", envFileRootPath, err)
	}
	if n <= 0 {
		return nil, fmt.Errorf("empty .env file")
	}

	env := string(buffer[:n])
	if env == "" {
		return nil, fmt.Errorf("empty .env file content")
	}
	scanner := bufio.NewScanner(strings.NewReader(env))
	if scanner == nil {
		return nil, fmt.Errorf("failed to create scanner")
	}

	envConf, err := readEnvFile(scanner)
	if err != nil {
		return nil, fmt.Errorf("failed to read env file: %s", err.Error())
	}
	if envConf.Port <= 0 {
		return nil, fmt.Errorf("PORT is required")
	}

	return envConf, nil
}

func readEnvFile(scanner *bufio.Scanner) (*Config, error) {
	envConf := Config{Smtp: SmtpConfig{}, RateLimit: RpsLimit{}, TestUsers: make([]*TestUser, 0)}
	for scanner.Scan() {
		line := scanner.Text()
		kv := strings.Split(line, "=")
		if len(kv) != 2 {
			log.Println("invalid config line, ", line)
			continue
		}
		switch kv[0] {
		case "PORT":
			envConf.Port = parseInt(kv[0], kv[1])
		case "SQLITE":
			envConf.Sqlite = kv[1]
		case "SMTP_USER":
			envConf.Smtp.User = kv[1]
		case "SMTP_PASS":
			envConf.Smtp.Pass = kv[1]
		case "SMTP_HOST":
			envConf.Smtp.Host = kv[1]
		case "SMTP_PORT":
			envConf.Smtp.Port = kv[1]
		case "CACHE_SIZE":
			envConf.CacheSize = parseInt(kv[0], kv[1])
		case "THROTTLE_TOTAL_RPS":
			envConf.RateLimit.TotalRPS = parseInt(kv[0], kv[1])
		case "THROTTLE_TOTAL_BURST":
			envConf.RateLimit.TotalBurst = parseInt(kv[0], kv[1])
		case "THROTTLE_USER_RPS":
			envConf.RateLimit.UserRPS = parseInt(kv[0], kv[1])
		case "THROTTLE_USER_BURST":
			envConf.RateLimit.UserBurst = parseInt(kv[0], kv[1])
		case "TEST_DATA_INSERT":
			envConf.TestDataInsert = kv[1] == "true"
		case "TEST_USER":
			testUser := parseTestUser(kv[1])
			envConf.TestUsers = append(envConf.TestUsers, testUser)
		case "LOG_STDOUT":
			envConf.Log.Stdout = kv[1] == "true"
		case "LOG_DIR":
			envConf.Log.Dir = kv[1]
		default:
			log.Printf("unknown env config [%s]\n", line)
		}
	}
	return &envConf, nil
}

func parseInt(key string, val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("failed to parse int[%s] as [%s]: %v", val, key, err))
	}
	return i
}

// TEST_USER="name:ABCDE;email:abcd@gmail.com;pass:123456"
func parseTestUser(rawUser string) *TestUser {
	log.Printf("...digesting raw test user data [%s]", rawUser)
	testUser := TestUser{}
	rawUser = strings.Trim(rawUser, "\"")
	props := strings.Split(rawUser, ";")
	for _, prop := range props {
		var err error
		testUser, err = assignProp(testUser, prop)
		if err != nil {
			log.Fatalf("failed to parse test user section [%s]: %v\n", prop, err)
		}
	}
	if testUser.Name == "" {
		log.Fatalf("name is required for test user")
	}
	if testUser.Email == "" {
		log.Fatalf("email is required for test user")
	}
	if testUser.Pass == "" {
		log.Fatalf("pass is required")
	}
	return &testUser
}

func assignProp(testUser TestUser, prop string) (TestUser, error) {
	kv := strings.Split(prop, ":")
	if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
		return testUser, fmt.Errorf("invalid test user, %s", testUser)
	}
	switch kv[0] {
	case "name":
		testUser.Name = kv[1]
	case "email":
		testUser.Email = kv[1]
	case "pass":
		testUser.Pass = kv[1]
	default:
		log.Printf("unknown test user section [%s]\n", testUser)
	}

	return testUser, nil
}
