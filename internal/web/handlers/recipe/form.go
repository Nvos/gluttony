package recipe

import (
	"fmt"
	"gluttony/internal/recipe"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Form struct {
	ID                int64
	Name              string
	Description       string
	Source            string
	Instructions      string
	ThumbnailImageURL string
	Servings          int8
	PreparationTime   time.Duration
	CookTime          time.Duration
	Tags              []string
	Ingredients       []recipe.Ingredient
	Nutrition         recipe.Nutrition
}

func (form Form) ToInput(ownerID int64) recipe.CreateInput {
	return recipe.CreateInput{
		Name:            form.Name,
		Description:     form.Description,
		Source:          form.Source,
		Instructions:    form.Instructions,
		Servings:        form.Servings,
		PreparationTime: form.PreparationTime,
		CookTime:        form.CookTime,
		Tags:            form.Tags,
		Ingredients:     form.Ingredients,
		Nutrition:       form.Nutrition,
		OwnerID:         ownerID,
	}
}

func NewRecipeForm(values url.Values) (Form, error) {
	ingredients := make([]recipe.Ingredient, len(values["ingredient"]))

	quantities := values["quantity"]
	notes := values["note"]
	units := values["unit"]
	for i, name := range values["ingredient"] {
		quantity, err := strconv.ParseFloat(quantities[i], 32)
		if err != nil {
			return Form{}, fmt.Errorf("parse quantity: %w", err)
		}

		ingredients[i].Order = int8(i)
		ingredients[i].Quantity = float32(quantity)
		ingredients[i].Unit = units[i]
		ingredients[i].Name = name
		ingredients[i].Note = notes[i]
	}

	servings, err := strconv.ParseInt(values.Get("servings"), 10, 8)
	if err != nil {
		return Form{}, fmt.Errorf("parse servings: %w", err)
	}

	// TODO: handle errors
	preparationDuration, _ := ParseFormDuration(values.Get("preparation-time"))
	cookDuration, _ := ParseFormDuration(values.Get("cook-time"))
	calories, _ := strconv.ParseFloat(values.Get("calories"), 32)
	protein, _ := strconv.ParseFloat(values.Get("protein"), 32)
	fat, _ := strconv.ParseFloat(values.Get("fat"), 32)
	carbs, _ := strconv.ParseFloat(values.Get("carbs"), 32)

	id, err := strconv.ParseInt(values.Get("id"), 10, 64)
	if err != nil {
		return Form{}, fmt.Errorf("parse id: %w", err)
	}

	return Form{
		ID:                id,
		Name:              values.Get("name"),
		Description:       values.Get("description"),
		Source:            values.Get("source"),
		Instructions:      values.Get("instructions"),
		Servings:          int8(servings),
		PreparationTime:   preparationDuration,
		CookTime:          cookDuration,
		Tags:              values["tag"],
		Ingredients:       ingredients,
		ThumbnailImageURL: values.Get("cover-image-url"),
		Nutrition: recipe.Nutrition{
			Calories: float32(calories),
			Fat:      float32(fat),
			Carbs:    float32(carbs),
			Protein:  float32(protein),
		},
	}, nil
}

// TODO: move to some time utils
func ParseFormDuration(value string) (time.Duration, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("could not parse form value: %s, expected 2 parts", value)
	}

	return time.ParseDuration(fmt.Sprintf("%sh%sm", parts[0], parts[1]))
}
