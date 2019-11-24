package dependency

import (
	"github.com/SananGuliyev/goddd/domain/repository"
	"github.com/SananGuliyev/goddd/infrastructure/persistence/postgres"
	"github.com/go-pg/pg/v9"
)

func NewRecipeRepository(db interface{}) repository.RecipeRepository {
	switch connection := db.(type) {
	case *pg.DB:
		return postgres.NewRecipeRepository(connection)
	}

	return nil
}

func NewRatingRepository(db interface{}) repository.RatingRepository {
	switch connection := db.(type) {
	case *pg.DB:
		return postgres.NewRatingRepository(connection)
	}

	return nil
}
