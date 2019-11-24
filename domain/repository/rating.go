package repository

import (
	"github.com/SananGuliyev/goddd/domain/entity"
)

type RatingRepository interface {
	GetRatingByRecipeId(recipeId string) float32
	Save(*entity.Rating) error
}
