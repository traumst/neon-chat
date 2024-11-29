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
	RateLimit      RpsLimit
	TestUsers      TestUsers
	TestDataInsert bool
	BackupConfig   BackupConfig
}

func (config *Config) String() string {
	acc := fmt.Sprintln("port:", config.Port)
	acc += fmt.Sprintln("dbfile:", config.Sqlite)
	acc += fmt.Sprintln("stdout:", config.Log.Stdout)
	acc += fmt.Sprintln("stdout:", config.Log.Dir)
	acc += fmt.Sprintln("cache:", config.CacheSize)
	acc += fmt.Sprintln("rateLimits:", config.RateLimit)
	acc += fmt.Sprintln("testUser:", config.TestUsers)
	acc += fmt.Sprintln("testDataInsert:", config.TestDataInsert)
	acc += fmt.Sprintln("backups:", config.BackupConfig)
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

type RpsLimit struct {
	TotalRPS   int
	TotalBurst int
	UserRPS    int
	UserBurst  int
}

func (l RpsLimit) String() string {
	return fmt.Sprintf("totalRPS:%d,totalBurst:%d,userRPS:%d,userBurst:%d", l.TotalRPS, l.TotalBurst, l.UserRPS, l.UserBurst)
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

type BackupConfig struct {
	SessionFilePath  string
	UserChatFilePath string
}

func (bc BackupConfig) String() string {
	return fmt.Sprintf("sessions:%s,userChats:%s", bc.SessionFilePath, bc.UserChatFilePath)
}
