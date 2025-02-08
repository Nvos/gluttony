package recipe

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"os"
	"path/filepath"
)

func NewSearchIndex(workDir string) (bleve.Index, error) {
	mapping := bleve.NewIndexMapping()
	var index bleve.Index

	indexPath := filepath.Join(workDir, "recipe-index.bleve")
	_, err := os.Stat(indexPath)
	if os.IsNotExist(err) {
		index, err = bleve.New(indexPath, mapping)
		if err != nil {
			return nil, fmt.Errorf("bleve create new: %w", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("bleve stat index file: %w", err)
	}

	if index == nil {
		index, err = bleve.Open(indexPath)
		if err != nil {
			return nil, fmt.Errorf("bleve open index file: %w", err)
		}
	}

	return index, nil
}
