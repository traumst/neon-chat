package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port      int
	LoadLocal bool
	Sqlite    string
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

func ArgsRead() (*Config, error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("no arguments provided")
	}

	argName := ""
	args := &Config{}
	err := error(nil)
	for _, arg := range os.Args[1:] {
		if err != nil {
			break
		}
		if strings.HasPrefix(arg, "-") {
			argName = arg
			continue
		}
		switch argName {
		case "-p", "--port":
			args.Port, err = strconv.Atoi(arg)
			if err == nil && args.Port > 0 {
				continue
			}
			err = fmt.Errorf("invalid --port value [%s], %v", arg, err)
		case "-h", "--help":
			err = fmt.Errorf("help requested")
		default:
			err = fmt.Errorf("unknown argument: %s", argName)
		}
	}

	if err != nil {
		return nil, err
	}

	return args, nil
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

	envConf := Config{}
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
		default:
			log.Printf("	unknown env [%s]=[%s]\n", kv[0], kv[1])
		}
	}

	if envConf.Port <= 0 {
		return nil, fmt.Errorf("PORT is required")
	}

	return &envConf, nil
}
