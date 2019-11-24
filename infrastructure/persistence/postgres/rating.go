package postgres

import (
	"github.com/SananGuliyev/goddd/domain/entity"
	"github.com/SananGuliyev/goddd/domain/repository"
	"github.com/go-pg/pg/v9"
)

type ratingRepository struct {
	db *pg.DB
}

func (r *ratingRepository) GetRatingByRecipeId(recipeId string) float32 {
	var rating float32
	err := r.db.Model((*entity.Rating)(nil)).
		ColumnExpr("ROUND(CAST(SUM(value) AS DECIMAL) / COUNT(id), 2) AS rating").
		Where("recipe_id = ?", recipeId).
		Select(&rating)

	if err != nil {
		return 0
	}

	return rating
}

func (r *ratingRepository) Save(rating *entity.Rating) error {
	err := r.db.Insert(rating)

	if err != nil {
		return err
	}

	return nil
}

func NewRatingRepository(db *pg.DB) repository.RatingRepository {
	return &ratingRepository{db: db}
}
