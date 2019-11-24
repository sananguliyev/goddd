package repository

import (
	"github.com/SananGuliyev/goddd/domain/entity"
)

type RecipeRepository interface {
	GetAll(page int, limit int, filter string) (*[]entity.Recipe, error)
	GetById(id string) (*entity.Recipe, error)
	Save(*entity.Recipe) error
	Update(*entity.Recipe) error
	Delete(*entity.Recipe) error
}
