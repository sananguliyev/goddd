package postgres

import (
	"github.com/SananGuliyev/goddd/domain/entity"
	"github.com/SananGuliyev/goddd/domain/repository"
	"github.com/SananGuliyev/goddd/domain/throwable"
	"github.com/go-pg/pg/v9"
)

type recipeRepository struct {
	db     *pg.DB
	parser *parser
}

func (r *recipeRepository) GetAll(page int, limit int, filter string) (*[]entity.Recipe, error) {
	recipes := &[]entity.Recipe{}
	offset := (page - 1) * limit

	sql, parserErr := r.parser.Parse(filter, limit, offset)

	if parserErr != nil {
		return nil, parserErr
	}

	query := "SELECT * FROM recipes " + sql

	_, repositoryErr := r.db.Query(recipes, query)

	if repositoryErr != nil {
		return nil, repositoryErr
	}

	return recipes, nil
}

func (r *recipeRepository) GetById(id string) (*entity.Recipe, error) {
	recipe := &entity.Recipe{Id: id}
	err := r.db.Select(recipe)
	if err != nil {
		switch e := err.(type) {
		case pg.Error:
			return nil, e
		default:
			return nil, throwable.NewNotFound("Recipe not found")
		}
	}
	return recipe, nil
}

func (r *recipeRepository) Save(recipe *entity.Recipe) error {
	err := r.db.Insert(recipe)

	return err
}

func (r *recipeRepository) Update(recipe *entity.Recipe) error {
	err := r.db.Update(recipe)

	return err
}

func (r *recipeRepository) Delete(recipe *entity.Recipe) error {
	err := r.db.Delete(recipe)

	return err
}

func NewRecipeRepository(db *pg.DB) repository.RecipeRepository {
	return &recipeRepository{
		db:     db,
		parser: NewParser(),
	}
}
