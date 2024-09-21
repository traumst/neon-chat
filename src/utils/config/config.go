package config

import "fmt"

type Config struct {
	Port           int
	Sqlite         string
	Smtp           SmtpConfig
	CacheSize      int
	TestUsers      TestUsers
	TestDataInsert bool
}

func (a *Config) String() string {
	return fmt.Sprintf("{Port:%d,Sqlite:%s}", a.Port, a.Sqlite)
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
	Salt  string
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
