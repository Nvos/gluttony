package recipe

import (
	"database/sql"
	"gluttony/internal/recipe/queries"
)

type Service struct {
	db *queries.Queries
}

func NewService(db *sql.DB) *Service {
	return &Service{db: queries.New(db)}
}

func (s *Service) Create() error {
	return nil
}
