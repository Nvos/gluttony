package recipe

import (
	"fmt"
	"gluttony/internal/share"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Ingredient struct {
	Order    int8
	Name     string
	Quantity float32
	Unit     string
}

type Nutrition struct {
	Calories float32
	Fat      float32
	Carbs    float32
	Protein  float32
}

type CreateForm struct {
	Name            string
	Description     string
	Source          string
	Instructions    string
	CoverImageURL   string
	Servings        int8
	PreparationTime time.Duration
	CookTime        time.Duration
	Tags            []string
	Ingredients     []Ingredient
	Nutrition       Nutrition
}

type CreateFormModel struct {
	*share.Context
	Form CreateForm
}

func RecipeCreateViewHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		appCtx := share.MustGetContext(r.Context())
		get, err := deps.templates.Get("recipe", "recipe_create")
		if err != nil {
			// TODO: proper err
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}

		model := CreateFormModel{
			Context: appCtx,
			Form: CreateForm{
				Ingredients: []Ingredient{
					{
						Name:     "",
						Quantity: 0,
						Unit:     "g",
					},
				},
			},
		}

		err = get.View(w, model)
		if err != nil {
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}
	}
}

func RecipesViewHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		appCtx := share.MustGetContext(r.Context())
		get, err := deps.templates.Get("recipe", "recipes")
		if err != nil {
			// TODO: proper err
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}

		err = get.View(w, appCtx)
		if err != nil {
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}
	}
}

func RecipesCreateHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << (10 * 2)); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		model := CreateFormModel{
			Context: share.MustGetContext(r.Context()),
			Form:    NewRecipeForm(r.MultipartForm.Value),
		}

		coverImage := r.MultipartForm.File["cover-image"]
		if len(coverImage) == 1 {
			file, err := coverImage[0].Open()
			if err != nil {
				// TODO: handle err
				panic(fmt.Errorf("could not open cover image: %v", err))
			}
			defer file.Close()
			coverImage[0].Header.Get("Content-Type")

			fileName, err := deps.mediaStore.Store(file)
			if err != nil {
				// TODO: handle err
				panic(fmt.Errorf("could not store cover image: %v", err))
			}

			model.Form.CoverImageURL = fmt.Sprintf("/media/%s", fileName)
		}

		get, err := deps.templates.Get("recipe", "recipe_create")
		if err != nil {
			// TODO: proper err
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}

		err = get.Fragment(w, "recipe-create/form", model)
		if err != nil {
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}
	}
}

func NewRecipeForm(values url.Values) CreateForm {
	ingredients := make([]Ingredient, len(values["ingredient"]))

	quantities := values["quantity"]
	units := values["unit"]
	for i, ingredient := range values["ingredient"] {
		value, err := strconv.ParseFloat(quantities[i], 32)
		if err != nil {
			// TODO: handle
			panic(err)
		}

		ingredients[i].Order = int8(i)
		ingredients[i].Quantity = float32(value)
		ingredients[i].Unit = units[i]
		ingredients[i].Name = ingredient
	}

	servings, err := strconv.ParseInt(values.Get("servings"), 10, 8)
	if err != nil {
		// TODO: handle
		panic(err)
	}

	// TODO: handle errors
	preparationDuration, _ := ParseFormDuration(values.Get("preparation-time"))
	cookDuration, _ := ParseFormDuration(values.Get("cook-time"))
	calories, _ := strconv.ParseFloat(values.Get("calories"), 32)
	protein, _ := strconv.ParseFloat(values.Get("protein"), 32)
	fat, _ := strconv.ParseFloat(values.Get("fat"), 32)
	carbs, _ := strconv.ParseFloat(values.Get("carbs"), 32)

	return CreateForm{
		Name:            values.Get("name"),
		Description:     values.Get("description"),
		Source:          values.Get("source"),
		Instructions:    values.Get("instructions"),
		Servings:        int8(servings),
		PreparationTime: preparationDuration,
		CookTime:        cookDuration,
		Tags:            values["tag"],
		Ingredients:     ingredients,
		CoverImageURL:   values.Get("cover-image-url"),
		Nutrition: Nutrition{
			Calories: float32(calories),
			Fat:      float32(fat),
			Carbs:    float32(carbs),
			Protein:  float32(protein),
		},
	}
}

// TODO: move to some time utils

func ParseFormDuration(value string) (time.Duration, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("could not parse form value: %s, expected 2 parts", value)
	}

	return time.ParseDuration(fmt.Sprintf("%sh%sm", parts[0], parts[1]))
}

type MediaStore interface {
	Store(file io.Reader) (string, error)
}
