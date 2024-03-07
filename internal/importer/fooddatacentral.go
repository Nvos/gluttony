package importer

import (
	"encoding/json"
	"fmt"
	"io"
)

type FoodCentralData struct {
	FoodClass   string `json:"foodClass"`
	Description string `json:"description"`
	// TODO(AK) 07/03/2024: more params?
}

func ImportFoodDataCentral(reader io.Reader) error {
	dec := json.NewDecoder(reader)
	_, err := dec.Token()
	if err != nil {
		return fmt.Errorf("consume open array bracket: %w", err)
	}

	var rows []FoodCentralData
	count := 0
	for dec.More() {
		var data FoodCentralData

		if err := dec.Decode(&data); err != nil {
			return fmt.Errorf("decode data: %w", err)
		}

		rows = append(rows, data)

		count++
	}

	_, err = dec.Token()
	if err != nil {
		return fmt.Errorf("consume close array bracker: %w", err)
	}

	return nil
}
