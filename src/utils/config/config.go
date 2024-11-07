package config

import (
	"fmt"
)

type Config struct {
	Log            LogConfig
	Port           int
	Sqlite         string
	Smtp           SmtpConfig
	CacheSize      int
	TestUsers      TestUsers
	TestDataInsert bool
}

func (config *Config) String() string {
	acc := fmt.Sprintln("port:", config.Port)
	acc += fmt.Sprintln("dbfile:", config.Sqlite)
	acc += fmt.Sprintln("stdout:", config.Log.Stdout)
	acc += fmt.Sprintln("stdout:", config.Log.Dir)
	acc += fmt.Sprintln("cache:", config.CacheSize)
	acc += fmt.Sprintln("testUser:", config.TestUsers)
	acc += fmt.Sprintln("testDataInsert:", config.TestDataInsert)
	return acc
}

type LogConfig struct {
	Stdout bool
	Dir    string
}

type SmtpConfig struct {
	User string
	Pass string
	Host string
	Port string
}

type TestUser struct {
	Name  string
	Email string
	Pass  string
}

type TestUsers []*TestUser

func (tu TestUsers) GetNames() []string {
	names := make([]string, 0)
	for _, u := range tu {
		names = append(names, u.Name)
	}
	return names
}

func (tu TestUsers) String() string {
	acc := "["
	for _, u := range tu {
		acc += fmt.Sprintf("\n name:%s,email:%s,pass:%s", u.Name, u.Email, u.Pass)
	}
	if len(tu) > 1 {
		acc += "\n"
	}
	acc += "]"
	return acc
}
