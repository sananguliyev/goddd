package postgres

import (
	"fmt"
	"github.com/SananGuliyev/goddd/config"
	"github.com/SananGuliyev/goddd/domain/entity"
	"github.com/SananGuliyev/goddd/domain/tool"
	tool3 "github.com/SananGuliyev/goddd/infrastructure/tool"
	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RecipeRepositoryTestSuite struct {
	suite.Suite
	db          *pg.DB
	parser      *parser
	idGenerator tool.IdGenerator
}

func (suite *RecipeRepositoryTestSuite) SetupTest() {
	host, port, user, password, database := config.GetDatabaseConfig()

	suite.db = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		User:     user,
		Password: password,
		Database: database,
	})
	suite.parser = NewParser()
	suite.idGenerator = tool3.NewUuidGenerator()
}

func (suite *RecipeRepositoryTestSuite) TestSaveReadDelete() {
	someId := suite.idGenerator.Generate()

	someRecipe := &entity.Recipe{
		Id:           someId,
		Name:         "someName",
		PrepareTime:  15,
		Difficulty:   15,
		IsVegetarian: true,
	}

	underTest := recipeRepository{db: suite.db}

	errSave := underTest.Save(someRecipe)
	recipe, errGet := underTest.GetById(someId)
	errDelete := underTest.Delete(recipe)
	recipeAfterDelete, errAfterDelete := underTest.GetById(someId)

	assert.Nil(suite.T(), errSave)
	assert.Nil(suite.T(), errGet)
	assert.Nil(suite.T(), errDelete)
	assert.Nil(suite.T(), recipeAfterDelete)
	assert.NotNil(suite.T(), errAfterDelete)
	assert.Equal(suite.T(), someRecipe, recipe)
}

func (suite *RecipeRepositoryTestSuite) TestSave_WithInvalidUuid() {
	recipe := &entity.Recipe{
		Id: "someInvalidId",
	}

	underTest := recipeRepository{db: suite.db}

	err := underTest.Save(recipe)

	assert.NotNil(suite.T(), err)
}

func TestRecipeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RecipeRepositoryTestSuite))
}
