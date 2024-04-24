package config

import (
	"gluttony/internal/x/validatex"
)

func ValidateConfig(cfg Config) error {
	return validatex.NewValidationError(
		// Database
		validatex.String("database.database", cfg.Database.Database, validatex.Empty()),
		validatex.String("database.host", cfg.Database.Host, validatex.Empty()),
		validatex.String("database.username", cfg.Database.User, validatex.Empty()),
		validatex.String("database.password", cfg.Database.Password, validatex.Empty()),
		validatex.String("database.port", cfg.Database.Port, validatex.Empty()),

		// Server
		validatex.String("server.port", cfg.Server.Port, validatex.Empty()),
	)
}
