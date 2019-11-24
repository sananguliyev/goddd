package interactor

import (
	"github.com/SananGuliyev/goddd/domain/entity"
	"github.com/SananGuliyev/goddd/domain/repository"
	"github.com/SananGuliyev/goddd/domain/tool"
)

type RecipeInteractor interface {
	Create(recipe *entity.Recipe) error
	Get(id string) (*entity.Recipe, float32, error)
	Update(id string, recipe *entity.Recipe) error
	Delete(id string) error
	List(page int, limit int, filter string) (*[]entity.Recipe, error)
	Rate(recipeId string, rating *entity.Rating) error
}

type recipeInteractor struct {
	recipeRepository repository.RecipeRepository
	ratingRepository repository.RatingRepository
	idGenerator      tool.IdGenerator
}

func (ri *recipeInteractor) Create(recipe *entity.Recipe) error {
	recipe.Id = ri.idGenerator.Generate()
	err := ri.recipeRepository.Save(recipe)

	if err != nil {
		return err
	}
	return nil
}

func (ri *recipeInteractor) Get(id string) (*entity.Recipe, float32, error) {
	recipe, err := ri.recipeRepository.GetById(id)

	if err != nil {
		return nil, 0, err
	}

	rating := ri.ratingRepository.GetRatingByRecipeId(id)

	return recipe, rating, nil
}

func (ri *recipeInteractor) Update(id string, recipe *entity.Recipe) error {
	var err error

	_, err = ri.recipeRepository.GetById(id)

	if err != nil {
		return err
	}

	recipe.Id = id

	err = ri.recipeRepository.Update(recipe)

	if err != nil {
		return err
	}

	return nil
}

func (ri *recipeInteractor) Delete(id string) error {
	var err error

	recipe, err := ri.recipeRepository.GetById(id)

	if err != nil {
		return err
	}

	err = ri.recipeRepository.Delete(recipe)
	if err != nil {
		return err
	}

	return nil
}

func (ri *recipeInteractor) List(page int, limit int, filter string) (*[]entity.Recipe, error) {
	recipes, err := ri.recipeRepository.GetAll(page, limit, filter)

	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (ri *recipeInteractor) Rate(recipeId string, rating *entity.Rating) error {
	var err error

	_, err = ri.recipeRepository.GetById(recipeId)

	if err != nil {
		return err
	}

	rating.Id = ri.idGenerator.Generate()
	rating.RecipeId = recipeId
	err = ri.ratingRepository.Save(rating)

	if err != nil {
		return err
	}
	return nil
}

func NewRecipeInteractor(
	recipeRepository repository.RecipeRepository,
	ratingRepository repository.RatingRepository,
	idGenerator tool.IdGenerator,
) RecipeInteractor {
	return &recipeInteractor{
		recipeRepository: recipeRepository,
		ratingRepository: ratingRepository,
		idGenerator:      idGenerator,
	}
}
