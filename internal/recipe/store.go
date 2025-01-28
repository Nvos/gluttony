package recipe

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
	"text/template"
)

var allRecipesQuery = template.Must(
	template.New("").
		Funcs(template.FuncMap{
			"array": func(count int) string {
				return strings.Repeat(",?", count)[1:]
			},
		}).
		Parse(`
		SELECT id, name, description, thumbnail_url
		FROM recipes
		{{if ne (len .RecipeIDs) 0}}
		WHERE id in ({{array (len .RecipeIDs)}})
		{{end}}
		ORDER BY id DESC
		LIMIT ? OFFSET ?;
	`),
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) AllRecipeSummaries(ctx context.Context, input SearchInput) ([]Summary, error) {
	var buffer bytes.Buffer
	if err := allRecipesQuery.Execute(&buffer, input); err != nil {
		return nil, err
	}

	query := buffer.String()
	params := make([]any, 0, 2+len(input.RecipeIDs))
	for i := range input.RecipeIDs {
		params = append(params, input.RecipeIDs[i])
	}
	params = append(params, input.Limit, input.Page*input.Limit)

	rows, err := s.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out := make([]Summary, 0, 20)
	for rows.Next() {
		value := Summary{}
		err = rows.Scan(
			&value.ID,
			&value.Name,
			&value.Description,
			&value.ThumbnailImageURL,
		)
		if err != nil {
			return nil, err
		}

		out = append(out, value)
	}

	return out, nil
}
