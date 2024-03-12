package config

import "gluttony/internal/util/validateutil"

func ValidateConfig(cfg Config) error {
	return validateutil.NewValidationError(
		// Database
		validateutil.String("database.database", cfg.Database.Database, validateutil.Empty()),
		validateutil.String("database.host", cfg.Database.Host, validateutil.Empty()),
		validateutil.String("database.username", cfg.Database.Username, validateutil.Empty()),
		validateutil.String("database.password", cfg.Database.Password, validateutil.Empty()),
		validateutil.Number("database.port", cfg.Database.Port, validateutil.Min(1, true)),

		// Server
		validateutil.Number("server.port", cfg.Server.Port, validateutil.Min(1, true)),
	)
}
