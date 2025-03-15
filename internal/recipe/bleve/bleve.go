package bleve

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"gluttony/internal/recipe"
	"gluttony/pkg/pagination"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var _ recipe.Index = (*Index)(nil)

type indexValue struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Index struct {
	index bleve.Index
}

func New(workDir string) (*Index, error) {
	var index bleve.Index

	indexPath := filepath.Join(workDir, "recipe-index.bleve")
	_, err := os.Stat(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			mapping := bleve.NewIndexMapping()
			index, err = bleve.New(indexPath, mapping)
			if err != nil {
				return nil, fmt.Errorf("bleve create new: %w", err)
			}

			return &Index{
				index: index,
			}, nil
		}

		return nil, fmt.Errorf("bleve stat index file: %w", err)
	}

	index, err = bleve.Open(indexPath)
	if err != nil {
		return nil, fmt.Errorf("bleve open index file: %w", err)
	}

	return &Index{index: index}, nil
}

func (idx *Index) Index(value recipe.Recipe) error {
	toIndex := indexValue{
		ID:          value.ID,
		Name:        value.Name,
		Description: value.Description,
	}

	if err := idx.index.Index(strconv.Itoa(int(value.ID)), toIndex); err != nil {
		return fmt.Errorf("recipe index: %w", err)
	}

	return nil
}

func (idx *Index) Close() error {
	if err := idx.index.Close(); err != nil {
		return fmt.Errorf("close recipe index: %w", err)
	}

	return nil
}

func (idx *Index) Search(
	ctx context.Context,
	phrase string,
	offset pagination.Offset,
) (recipe.SearchResult, error) {
	request := bleve.NewSearchRequestOptions(
		buildQuery(phrase),
		int(offset.Limit),
		int(offset.Offset),
		false,
	)

	searchResult, err := idx.index.SearchInContext(ctx, request)
	if err != nil {
		return recipe.SearchResult{}, fmt.Errorf("recipe search: %w", err)
	}

	out := recipe.SearchResult{
		TotalCount: int64(searchResult.Total),
		IDs:        make([]int32, 0, len(searchResult.Hits)),
	}
	for i := range searchResult.Hits {
		id, err := strconv.ParseInt(searchResult.Hits[i].ID, 10, 32)
		if err != nil {
			panic(fmt.Sprintf("unexpected recipe index id: %v", err))
		}

		out.IDs = append(out.IDs, int32(id))
	}
	bleve.NewDisjunctionQuery()
	return out, nil
}

func buildQuery(phrase string) *query.DisjunctionQuery {
	terms := strings.Split(phrase, " ")

	dj := query.NewDisjunctionQuery(nil)
	for i := range terms {
		nameMatch := bleve.NewMatchQuery(terms[i])
		nameMatch.SetField("name")
		nameMatch.SetBoost(3)

		descriptionMatch := bleve.NewMatchQuery(terms[i])
		descriptionMatch.SetField("description")

		dj.AddQuery(nameMatch, descriptionMatch)
	}

	return dj
}
