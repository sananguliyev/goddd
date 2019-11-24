package main

import (
	"fmt"
	"github.com/SananGuliyev/goddd/application/handler"
	"github.com/SananGuliyev/goddd/dependency"
	"github.com/SananGuliyev/goddd/domain/interactor"
	"github.com/SananGuliyev/goddd/presenter"
)

func main() {
	db, err := dependency.NewPostgresConnection()

	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	} else {
		defer dependency.Close(db)
	}

	recipeInteractor := interactor.NewRecipeInteractor(
		dependency.NewRecipeRepository(db),
		dependency.NewRatingRepository(db),
		dependency.NewIdGenerator(),
	)
	recipeHandler := handler.RecipeHandler{Interactor: recipeInteractor}

	httpServer := presenter.NewHttpServer(recipeHandler)
	httpServer.Run()
}
