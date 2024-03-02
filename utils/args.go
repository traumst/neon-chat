package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Args struct {
	Port int
}

func ArgsHelp() string {
	return `Usage: go.chat [options]
Options:
  -h / --help	- Show this help message

  -p / --port	- Port to listen on
  `
}

func ArgsRead() (*Args, error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("no arguments provided")
	}

	argName := ""
	args := &Args{}
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
