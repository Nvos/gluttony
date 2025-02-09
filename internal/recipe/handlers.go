package recipe

import (
	"fmt"
	"gluttony/internal/httputil"
	"gluttony/internal/ingredient"
	"gluttony/internal/share"
	"net/http"
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
	Ingredients       []Ingredient
	Nutrition         Nutrition
}

func (form Form) ToInput() CreateInput {
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
	Form Form
}

type UpdateModel struct {
	*share.Context
	Form Form
}

type ListModel struct {
	*share.Context
	Recipes []Summary
}

type ViewModel struct {
	*share.Context
	Recipe Full
}

func CreateViewHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		model := CreateModel{
			Context: share.MustGetContext(r.Context()),
			Form: Form{
				Servings: 1,
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

		return deps.templates.View(w, "recipe_create", model)
	}
}

func EditViewHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		recipeIDPathParam := r.PathValue("recipe_id")
		recipeID, err := strconv.ParseInt(recipeIDPathParam, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		recipe, err := deps.service.GetFull(r.Context(), recipeID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return nil
		}

		tags := make([]string, 0, len(recipe.Tags))
		for _, tag := range recipe.Tags {
			tags = append(tags, tag.Name)
		}

		model := UpdateModel{
			Context: share.MustGetContext(r.Context()),
			Form: Form{
				ID:                int64(recipe.ID),
				Name:              recipe.Name,
				Description:       recipe.Description,
				Source:            recipe.Source,
				Instructions:      recipe.Instructions,
				ThumbnailImageURL: recipe.ThumbnailImageURL,
				Servings:          recipe.Servings,
				PreparationTime:   recipe.PreparationTime,
				CookTime:          recipe.CookTime,
				Tags:              tags,
				Ingredients:       recipe.Ingredients,
				Nutrition:         recipe.Nutrition,
			},
		}

		return deps.templates.View(w, "recipe_edit", model)
	}
}

func RecipesViewHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		search := r.URL.Query().Get("search")
		pageParam := r.URL.Query().Get("page")
		limit := 20

		page := 0
		if pageParam != "" {
			pageInt, err := strconv.Atoi(pageParam)
			if err != nil {
				return fmt.Errorf("parse page to int: %w", err)
			}
			page = pageInt
		}

		recipePartials, err := deps.service.AllSummaries(r.Context(), SearchInput{
			Search: search,
			Page:   int64(page),
			Limit:  int64(limit),
		})
		if err != nil {
			return fmt.Errorf("could not get recipe partials: %w", err)
		}

		model := ListModel{
			Context: share.MustGetContext(r.Context()),
			Recipes: recipePartials,
		}

		if httputil.IsHTMXRequest(r) {
			return deps.templates.Fragment(w, "recipes/list", model)
		}

		return deps.templates.View(w, "recipes", model)
	}
}

func CreateFormHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		// TODO: extract max memory to const
		if err := r.ParseMultipartForm(1 << (10 * 2)); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		form, err := NewRecipeForm(r.MultipartForm.Value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		model := CreateModel{
			Context: share.MustGetContext(r.Context()),
			Form:    form,
		}

		input := model.Form.ToInput()

		coverImage := r.MultipartForm.File["thumbnail-image"]
		if len(coverImage) == 1 {
			file, err := coverImage[0].Open()
			if err != nil {
				// TODO: handle err
				panic(fmt.Errorf("could not open cover image: %w", err))
			}
			defer file.Close()

			input.ThumbnailImage = file
		}

		err = deps.service.Create(r.Context(), input)
		if err == nil {
			httputil.HTMXRedirect(w, "/recipes")
			return nil
		}

		return deps.templates.Fragment(w, "recipe/form", model)
	}
}

func UpdateFormHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		// TODO: extract max memory to const
		if err := r.ParseMultipartForm(1 << (10 * 2)); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		form, err := NewRecipeForm(r.MultipartForm.Value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		model := UpdateModel{
			Context: share.MustGetContext(r.Context()),
			Form:    form,
		}

		input := model.Form.ToInput()

		coverImage := r.MultipartForm.File["thumbnail-image"]
		if len(coverImage) == 1 {
			file, err := coverImage[0].Open()
			if err != nil {
				// TODO: handle err
				panic(fmt.Errorf("could not open cover image: %w", err))
			}
			defer file.Close()

			input.ThumbnailImage = file
		}

		err = deps.service.Update(r.Context(), UpdateInput{
			ID:          form.ID,
			CreateInput: input,
		})
		if err == nil {
			httputil.HTMXRedirect(w, fmt.Sprintf("/recipes/%d", form.ID))
			return nil
		}

		return deps.templates.Fragment(w, "recipe/form", model)
	}
}

func ViewHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		recipeIDRaw := r.PathValue("recipe_id")
		recipeID, err := strconv.Atoi(recipeIDRaw)
		if err != nil {
			// TODO: likely should return template or 404
			return fmt.Errorf("could not parse recipe id: %w", err)
		}

		recipe, err := deps.service.GetFull(r.Context(), int64(recipeID))
		if err != nil {
			// TODO: check for 404 (will need to add view + redirect?)
			return fmt.Errorf("could not get recipe partials: %w", err)
		}

		instructionsHTML, err := deps.markdown.ConvertToHTML(recipe.Instructions)
		if err != nil {
			return fmt.Errorf("could not convert instructions to HTML: %w", err)
		}

		recipe.Instructions = instructionsHTML
		model := ViewModel{
			Context: share.MustGetContext(r.Context()),
			Recipe:  recipe,
		}

		return deps.templates.View(w, "recipe", model)
	}
}

func NewRecipeForm(values url.Values) (Form, error) {
	ingredients := make([]Ingredient, len(values["ingredient"]))

	quantities := values["quantity"]
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
		Nutrition: Nutrition{
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
