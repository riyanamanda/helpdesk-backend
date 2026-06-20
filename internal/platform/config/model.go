package config

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	App         App
	Database    Database
	IhsDatabase Database
	Auth        Auth
	Storage     Storage
	Redis       Redis
	Email       Email
	RabbitMQ    RabbitMQ
	Antrol      Antrol
}

type Antrol struct {
	Domain   string
	Username string
	Password string
}


type App struct {
	Name        string
	Host        string
	Port        string
	CORSOrigins []string
}

type Database struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
	SSLMode  string
}

type Auth struct {
	JWTSecret               string
	JWTExp                  time.Duration
	FirebaseProjectID       string
	FirebaseCredentialsJSON string
}

type Storage struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type Redis struct {
	Host     string
	Port     string
	Password string
}

type Email struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	UseSSL   bool
}

type RabbitMQ struct {
	Host     string
	Port     string
	Username string
	Password string
	VHost    string
}

func (r RabbitMQ) DSN() string {
	vhost := strings.TrimPrefix(r.VHost, "/")
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s", r.Username, r.Password, r.Host, r.Port, vhost)
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

func (d Database) MySqlConnString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
}
