package recipe

import "gluttony/internal/recipe/queries"

func NewSummaryFromDBRow(row queries.AllRecipeSummaryRow) Summary {
	out := Summary{
		ID:          int(row.ID),
		Name:        row.Name,
		Description: row.Description,
	}

	if row.ThumbnailUrl.Valid {
		out.ThumbnailImageURL = row.ThumbnailUrl.String
	}

	return out
}
