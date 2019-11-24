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

type RatingRepositoryTestSuite struct {
	suite.Suite
	db          *pg.DB
	idGenerator tool.IdGenerator
}

func (suite *RatingRepositoryTestSuite) SetupTest() {
	host, port, user, password, database := config.GetDatabaseConfig()

	suite.db = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		User:     user,
		Password: password,
		Database: database,
	})
	suite.idGenerator = tool3.NewUuidGenerator()
}

func (suite *RatingRepositoryTestSuite) TestGetRatingByRecipeId() {
	someRecipeId := suite.idGenerator.Generate()
	someFirstValue := int8(5)
	someSecondValue := int8(4)
	expected := float32(someFirstValue+someSecondValue) / float32(2)

	underTest := ratingRepository{db: suite.db}

	err1 := underTest.Save(&entity.Rating{
		Id:       suite.idGenerator.Generate(),
		RecipeId: someRecipeId,
		Value:    someFirstValue,
	})

	err2 := underTest.Save(&entity.Rating{
		Id:       suite.idGenerator.Generate(),
		RecipeId: someRecipeId,
		Value:    someSecondValue,
	})

	actual := underTest.GetRatingByRecipeId(someRecipeId)

	assert.Nil(suite.T(), err1)
	assert.Nil(suite.T(), err2)
	assert.Equal(suite.T(), expected, actual)
}

func (suite *RatingRepositoryTestSuite) TestSave() {
	rating := &entity.Rating{
		Id:       suite.idGenerator.Generate(),
		RecipeId: suite.idGenerator.Generate(),
		Value:    5,
	}

	underTest := ratingRepository{db: suite.db}

	err := underTest.Save(rating)

	assert.Nil(suite.T(), err)
}

func (suite *RatingRepositoryTestSuite) TestSave_WithInvalidUuid() {
	rating := &entity.Rating{
		Id:       suite.idGenerator.Generate(),
		RecipeId: "someInvalidUuid",
		Value:    5,
	}

	underTest := ratingRepository{db: suite.db}

	err := underTest.Save(rating)

	assert.NotNil(suite.T(), err)
}

func TestRatingRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RatingRepositoryTestSuite))
}
