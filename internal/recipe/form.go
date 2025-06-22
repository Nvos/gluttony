package recipe

import (
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Form struct {
	ID              int32
	Name            string
	Description     string
	Source          string
	Instructions    string
	ThumbnailImage  *multipart.FileHeader
	Servings        int8
	PreparationTime time.Duration
	CookTime        time.Duration
	Tags            []string
	Ingredients     []Ingredient
	Nutrition       Nutrition
}

func (form Form) ToInput(ownerID int32) CreateInput {
	return CreateInput{
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
		ThumbnailImage:  form.ThumbnailImage,
	}
}

func NewRecipeForm(form *multipart.Form) (Form, error) {
	var values url.Values = form.Value
	ingredients := make([]Ingredient, len(values["ingredient"]))

	quantities := values["quantity"]
	notes := values["note"]
	units := values["unit"]
	for i, name := range values["ingredient"] {
		quantity, err := strconv.ParseFloat(quantities[i], 32)
		if err != nil {
			return Form{}, fmt.Errorf("parse quantity: %w", err)
		}

		if i > math.MaxInt8 {
			return Form{}, errors.New("too many ingredients")
		}

		ingredients[i].Order = int8(i) // #nosec G115 -- bounds checked above
		ingredients[i].Quantity = float32(quantity)
		ingredients[i].Unit = units[i]
		ingredients[i].Name = name
		ingredients[i].Note = notes[i]
	}

	servings, err := strconv.ParseInt(values.Get("servings"), 10, 8)
	if err != nil {
		return Form{}, fmt.Errorf("parse servings: %w", err)
	}

	preparationDuration, err := ParseFormDuration(values.Get("preparation-time"))
	if err != nil {
		return Form{}, fmt.Errorf("parse preparation duration: %w", err)
	}

	cookDuration, err := ParseFormDuration(values.Get("cook-time"))
	if err != nil {
		return Form{}, fmt.Errorf("parse cook duration: %w", err)
	}

	calories, err := strconv.ParseFloat(values.Get("calories"), 32)
	if err != nil {
		return Form{}, fmt.Errorf("parse calories: %w", err)
	}

	protein, err := strconv.ParseFloat(values.Get("protein"), 32)
	if err != nil {
		return Form{}, fmt.Errorf("parse protein: %w", err)
	}

	fat, err := strconv.ParseFloat(values.Get("fat"), 32)
	if err != nil {
		return Form{}, fmt.Errorf("parse fat: %w", err)
	}

	carbs, err := strconv.ParseFloat(values.Get("carbs"), 32)
	if err != nil {
		return Form{}, fmt.Errorf("parse carbs: %w", err)
	}

	id, err := strconv.ParseInt(values.Get("id"), 10, 32)
	if err != nil {
		return Form{}, fmt.Errorf("parse id: %w", err)
	}

	return Form{
		ID:              int32(id),
		Name:            values.Get("name"),
		Description:     values.Get("description"),
		Source:          values.Get("source"),
		Instructions:    values.Get("instructions"),
		Servings:        int8(servings),
		PreparationTime: preparationDuration,
		CookTime:        cookDuration,
		Tags:            values["tag"],
		Ingredients:     ingredients,
		ThumbnailImage:  GetThumbnail(form),
		Nutrition: Nutrition{
			Calories: float32(calories),
			Fat:      float32(fat),
			Carbs:    float32(carbs),
			Protein:  float32(protein),
		},
	}, nil
}

func ParseFormDuration(value string) (time.Duration, error) {
	const expectedPartCount = 2

	parts := strings.Split(value, ":")
	if len(parts) != expectedPartCount {
		return 0, fmt.Errorf("could not parse form value: %s, expected 2 parts", value)
	}

	duration, err := time.ParseDuration(fmt.Sprintf("%sh%sm", parts[0], parts[1]))
	if err != nil {
		return 0, fmt.Errorf("could not parse form duration: %w", err)
	}

	return duration, nil
}

func GetThumbnail(form *multipart.Form) *multipart.FileHeader {
	coverImage := form.File["thumbnail-image"]
	if len(coverImage) == 1 {
		return coverImage[0]
	}

	return nil
}