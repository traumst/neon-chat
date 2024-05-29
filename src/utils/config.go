package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type SmtpConfig struct {
	User string
	Pass string
	Host string
	Port string
}

type Config struct {
	Port      int
	LoadLocal bool
	Sqlite    string
	Smtp      SmtpConfig
}

func (a *Config) String() string {
	return fmt.Sprintf("{LoadLocal:%t,Port:%d,Sqlite:%s}", a.LoadLocal, a.Port, a.Sqlite)
}

func Help() string {
	return `By default, the application will read the config from the .env file in the root directory. 
	To set them:
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

	envConf := Config{Smtp: SmtpConfig{}}
	for scanner.Scan() {
		line := scanner.Text()
		kv := strings.Split(line, "=")
		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "PORT":
			port, err := strconv.Atoi(kv[1])
			if err != nil {
				return nil, fmt.Errorf("invalid PORT value [%s], %v", kv[1], err)
			}
			envConf.Port = port
		case "LOCAL":
			envConf.LoadLocal = strings.ToLower(kv[1]) == "true"
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
		default:
			log.Printf("unknown env config [%s]\n", line)
		}
	}

	if envConf.Port <= 0 {
		return nil, fmt.Errorf("PORT is required")
	}

	return &envConf, nil
}
