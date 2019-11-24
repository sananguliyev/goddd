package presenter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SananGuliyev/goddd/application/handler"
	"github.com/SananGuliyev/goddd/dependency"
	"github.com/SananGuliyev/goddd/domain/entity"
	"github.com/SananGuliyev/goddd/domain/interactor"
	"github.com/SananGuliyev/goddd/domain/repository"
	"github.com/SananGuliyev/goddd/domain/tool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MuxTestSuite struct {
	suite.Suite
	recipeRepository repository.RecipeRepository
	ratingRepository repository.RatingRepository
	idGenerator      tool.IdGenerator
	recipeHandler    handler.RecipeHandler
}

func (suite *MuxTestSuite) SetupTest() {
	db, _ := dependency.NewPostgresConnection()

	suite.recipeRepository = dependency.NewRecipeRepository(db)
	suite.ratingRepository = dependency.NewRatingRepository(db)
	suite.idGenerator = dependency.NewIdGenerator()
	recipeInteractor := interactor.NewRecipeInteractor(
		suite.recipeRepository,
		suite.ratingRepository,
		suite.idGenerator,
	)
	suite.recipeHandler = handler.RecipeHandler{Interactor: recipeInteractor}

}

func (suite *MuxTestSuite) TestRouteHealthz() {
	server := httptest.NewServer(getRouter(suite.recipeHandler))
	defer server.Close()

	response, err := server.Client().Get(fmt.Sprintf("%s/healthz", server.URL))

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, response.StatusCode)

	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(suite.T(), "OK", string(body))
}

func (suite *MuxTestSuite) TestRouteCreate() {
	expected := entity.Recipe{
		Name:         "GoDDD",
		PrepareTime:  15,
		Difficulty:   5,
		IsVegetarian: true,
	}
	recipeJson, _ := json.Marshal(expected)
	body := bytes.NewBuffer(recipeJson)

	server := httptest.NewServer(getRouter(suite.recipeHandler))
	defer server.Close()

	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/recipes", server.URL), body)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Authorization", "GoDDD")
	response, err := server.Client().Do(request)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, response.StatusCode)

	actual := &struct {
		Status int
		Data   entity.Recipe
	}{}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(actual)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), expected.Name, actual.Data.Name)
	assert.NotEmpty(suite.T(), actual.Data.Id)
}

func TestMuxTestSuite(t *testing.T) {
	suite.Run(t, new(MuxTestSuite))
}
