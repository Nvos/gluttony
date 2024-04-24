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
	Host string
	Port string
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Options  string
}
