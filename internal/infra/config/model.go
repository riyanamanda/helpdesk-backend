package config

import "fmt"

type Config struct {
	App      App
	Database Database
}

type App struct {
	Name string
	Host string
	Port string
}

type Database struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
	SSLMode  string
}

func (d Database) ConnString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host,
		d.Port,
		d.Username,
		d.Password,
		d.Name,
		d.SSLMode,
	)
}
