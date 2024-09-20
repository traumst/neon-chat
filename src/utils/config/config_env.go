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
	envConf := Config{Smtp: SmtpConfig{}, TestUsers: make([]*TestUser, 0)}
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
		case "CACHE_SIZE":
			envConf.CacheSize = parseInt(kv[0], kv[1])
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
		case "TEST_USER":
			envConf.TestUsers = accTestUsers(envConf.TestUsers, kv[1])
		default:
			log.Printf("unknown env config [%s]\n", line)
		}
	}
	return &envConf, nil
}

// expected format: val="user:ABCDE;email:abcd@gmail.com;pass:123456"
func accTestUsers(acc TestUsers, val string) TestUsers {
	sections := strings.Split(val, ";")
	if acc == nil {
		acc = make([]*TestUser, 0)
	}
	for _, section := range sections {
		user, err := parseTestUser(section)
		if err != nil {
			log.Printf("failed to parse test user section [%s]: %v\n", section, err)
			continue
		}
		acc = append(acc, user)
	}
	return acc
}

func parseTestUser(val string) (*TestUser, error) {
	kv := strings.Split(val, ":")
	if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
		return nil, fmt.Errorf("invalid test user, %s", val)
	}
	var user TestUser
	switch kv[0] {
	case "user":
		user = TestUser{Name: kv[1]}
	case "email":
		user.Email = kv[1]
	case "pass":
		user.Pass = kv[1]
	case "salt":
		user.Salt = kv[1]
	default:
		log.Printf("unknown test user section [%s]\n", val)
	}
	return &user, nil
}

func parseInt(key string, val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("failed to parse int[%s] as [%s]: %v", val, key, err))
	}
	return i
}
