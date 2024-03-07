package importer

import (
	"os"
	"testing"
)

func TestImportFoodDataCentral(t *testing.T) {
	open, err := os.Open("/home/czort/Downloads/FoodData_Central_foundation_food_json_2023-10-26/foundationDownload.json")
	if err != nil {
		t.Errorf("Invalid file: %v", err)
	}

	t.Cleanup(func() {
		_ = open.Close()
	})

	if err := ImportFoodDataCentral(open); err != nil {
		t.Errorf("unexpected ImportFoodDataCentral() error = %v", err)
	}
}
