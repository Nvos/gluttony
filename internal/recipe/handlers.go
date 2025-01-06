package recipe

import (
	"fmt"
	"gluttony/internal/ingredient"
	"gluttony/internal/share"
	"gluttony/x/httpx"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Ingredient struct {
	ingredient.Ingredient

	Order    int8
	Quantity float32
	// TODO: unit enum
	Unit string
}

type Nutrition struct {
	Calories float32
	Fat      float32
	Carbs    float32
	Protein  float32
}

type CreateForm struct {
	Name              string
	Description       string
	Source            string
	Instructions      string
	ThumbnailImageURL string
	Servings          int8
	PreparationTime   time.Duration
	CookTime          time.Duration
	Tags              []string
	Ingredients       []Ingredient
	Nutrition         Nutrition
}

func (form CreateForm) ToInput() CreateInput {
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
	}
}

type CreateModel struct {
	*share.Context
	Form CreateForm
}

type ListModel struct {
	*share.Context
	Recipes []Partial
}

func RecipeCreateViewHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		appCtx := share.MustGetContext(r.Context())
		get, err := deps.templates.Get("recipe", "recipe_create")
		if err != nil {
			// TODO: proper err
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}

		model := CreateModel{
			Context: appCtx,
			Form: CreateForm{
				Ingredients: []Ingredient{
					{
						Ingredient: ingredient.Ingredient{
							Name: "",
						},
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
		get, err := deps.templates.Get("recipe", "recipes")
		if err != nil {
			// TODO: proper err
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}

		recipePartials, err := deps.service.AllPartial(r.Context(), SearchInput{
			Query: r.URL.Query().Get("query"),
		})
		if err != nil {
			// TODO: proper err
			panic(fmt.Errorf("could not get recipe partials: %v", err))
		}

		model := ListModel{
			Context: share.MustGetContext(r.Context()),
			Recipes: recipePartials,
		}

		if httpx.IsHTMXRequest(r) {
			if err := get.Fragment(w, "recipes/list", model); err != nil {
				// TODO: proper err
				panic(fmt.Errorf("could not list recipes: %v", err))
			}

			return
		}

		err = get.View(w, model)
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

		model := CreateModel{
			Context: share.MustGetContext(r.Context()),
			Form:    NewRecipeForm(r.MultipartForm.Value),
		}

		input := model.Form.ToInput()

		coverImage := r.MultipartForm.File["thumbnail-image"]
		if len(coverImage) == 1 {
			file, err := coverImage[0].Open()
			if err != nil {
				// TODO: handle err
				panic(fmt.Errorf("could not open cover image: %v", err))
			}
			defer file.Close()

			input.ThumbnailImage = file
		}

		err := deps.service.Create(r.Context(), input)
		if err == nil {
			httpx.HTMXRedirect(w, "/recipes")

			return
		}
		println(fmt.Sprintf("could not create recipe template: %v", err))

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
	for i, name := range values["ingredient"] {
		quantity, err := strconv.ParseFloat(quantities[i], 32)
		if err != nil {
			// TODO: handle
			panic(err)
		}

		ingredients[i].Order = int8(i)
		ingredients[i].Quantity = float32(quantity)
		ingredients[i].Unit = units[i]
		ingredients[i].Name = name
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
	UploadImage(file io.Reader) (string, error)
}
