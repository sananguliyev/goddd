package handler

import (
	"encoding/json"
	"github.com/SananGuliyev/goddd/domain/entity"
	"github.com/SananGuliyev/goddd/domain/interactor"
	"net/http"
	"strconv"
)

type RecipeHandler struct {
	Handler
	Interactor interactor.RecipeInteractor
}

func (h RecipeHandler) Get(writer http.ResponseWriter, request *http.Request, id string) {
	recipe, rating, err := h.Interactor.Get(id)
	if err != nil {
		h.Error(writer, err)
		return
	}

	response := struct {
		Recipe *entity.Recipe `json:"recipe"`
		Rating float32        `json:"rating"`
	}{Recipe: recipe, Rating: rating}

	h.Data(writer, http.StatusOK, response)
}

func (h RecipeHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var err error
	recipe := &entity.Recipe{}

	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(recipe)

	if err != nil {
		h.Error(writer, err)
		return
	}

	if recipe.IsVegetarian == nil {
		isVegetarian := false
		recipe.IsVegetarian = &isVegetarian
	}

	err = h.Interactor.Create(recipe)
	if err != nil {
		h.Error(writer, err)
		return
	}

	h.Data(writer, http.StatusCreated, recipe)
}

func (h RecipeHandler) Update(writer http.ResponseWriter, request *http.Request, id string) {
	var err error
	recipe := &entity.Recipe{}

	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(recipe)

	if err != nil {
		h.Error(writer, err)
		return
	}

	err = h.Interactor.Update(id, recipe)
	if err != nil {
		h.Error(writer, err)
		return
	}

	h.Data(writer, http.StatusOK, recipe)
}

func (h RecipeHandler) Delete(writer http.ResponseWriter, request *http.Request, id string) {
	err := h.Interactor.Delete(id)
	if err != nil {
		h.Error(writer, err)
		return
	}

	h.Data(writer, http.StatusOK, struct{}{})
}

func (h RecipeHandler) List(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	page, invalidPage := strconv.Atoi(query.Get("page"))
	limit, invalidLimit := strconv.Atoi(query.Get("limit"))
	filter := query.Get("filter")

	if invalidPage != nil || page < 1 {
		page = 1
	}

	if invalidLimit != nil || limit < 1 {
		limit = 10 // Should be constant or come from configuration
	}

	recipes, err := h.Interactor.List(page, limit, filter)
	if err != nil {
		h.Error(writer, err)
		return
	}

	h.Data(writer, http.StatusOK, recipes)
}

func (h RecipeHandler) Rate(writer http.ResponseWriter, request *http.Request, recipeId string) {
	var err error

	rating := &entity.Rating{}
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(rating)

	if err != nil {
		h.Error(writer, err)
		return
	}

	err = h.Interactor.Rate(recipeId, rating)

	if err != nil {
		h.Error(writer, err)
		return
	}

	h.Data(writer, http.StatusOK, struct{}{})
}
