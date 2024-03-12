package config

import "log/slog"

type Config struct {
	Database Database
	Server   Server
	Logger   Logger
}

type Logger struct {
	Level slog.Level
}

type Server struct {
	Port int
}

type Database struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}
