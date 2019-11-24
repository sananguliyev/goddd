package interactor

import (
	"errors"
	"github.com/SananGuliyev/goddd/domain/entity"
	"github.com/SananGuliyev/goddd/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRecipeInteractor(t *testing.T) {
	expectedRecipeRepository := &mocks.RecipeRepository{}
	expectedIdGenerator := &mocks.IdGenerator{}
	expectedRatingRepository := &mocks.RatingRepository{}

	expected := &recipeInteractor{expectedRecipeRepository, expectedRatingRepository, expectedIdGenerator}

	actual := NewRecipeInteractor(expectedRecipeRepository, expectedRatingRepository, expectedIdGenerator)

	assert.Equal(t, expected, actual)
}

func TestRecipeInteractor_Create(t *testing.T) {
	var tests = []struct {
		name        string
		expectedErr error
	}{
		{"success case", nil},
		{"error case", errors.New("some error")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedId := "someId"
			someRecipe := &entity.Recipe{}

			mockIdGenerator := &mocks.IdGenerator{}
			mockIdGenerator.On("Generate").Return(expectedId).Once()

			mockRecipeRepository := &mocks.RecipeRepository{}
			mockRecipeRepository.On("Save", someRecipe).Return(test.expectedErr).Once()
			underTest := &recipeInteractor{recipeRepository: mockRecipeRepository, idGenerator: mockIdGenerator}
			actualErr := underTest.Create(someRecipe)
			actualId := someRecipe.Id

			assert.Equal(t, expectedId, actualId)
			assert.Equal(t, test.expectedErr, actualErr)
			mockIdGenerator.AssertExpectations(t)
			mockRecipeRepository.AssertExpectations(t)
		})
	}
}

func TestRecipeInteractor_Get(t *testing.T) {
	var tests = []struct {
		name           string
		someId         string
		expectedRecipe *entity.Recipe
		expectedRating float32
		expectedErr    error
	}{
		{
			name:           "found",
			someId:         "someExistingId",
			expectedRecipe: &entity.Recipe{Id: "someExistingId"},
			expectedRating: float32(4),
			expectedErr:    nil,
		},
		{
			name:           "not found",
			someId:         "someNonExistingId",
			expectedRecipe: nil,
			expectedRating: 0,
			expectedErr:    errors.New("some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRecipeRepository := &mocks.RecipeRepository{}
			mockRecipeRepository.On("GetById", test.someId).
				Return(test.expectedRecipe, test.expectedErr).
				Once()

			mockRatingRepository := &mocks.RatingRepository{}
			mockRatingRepository.On("GetRatingByRecipeId", test.someId).
				Return(test.expectedRating).
				Once()

			underTest := &recipeInteractor{
				recipeRepository: mockRecipeRepository,
				ratingRepository: mockRatingRepository,
			}

			actualRecipe, actualRating, actualErr := underTest.Get(test.someId)

			assert.Equal(t, test.expectedRecipe, actualRecipe)
			assert.Equal(t, test.expectedRating, actualRating)
			assert.Equal(t, test.expectedErr, actualErr)
			mockRecipeRepository.AssertExpectations(t)
			if test.expectedErr != nil {
				mockRatingRepository.AssertNotCalled(t, "GetRatingByRecipeId")
			} else {
				mockRatingRepository.AssertExpectations(t)
			}
		})
	}
}

func TestRecipeInteractor_Update(t *testing.T) {
	var tests = []struct {
		name        string
		expectedId  string
		expectedErr error
	}{
		{
			name:        "updated",
			expectedId:  "someExistingId",
			expectedErr: nil,
		},
		{
			name:        "failed",
			expectedId:  "someExistingId",
			expectedErr: errors.New("not updated"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			someRecipe := &entity.Recipe{}

			mockRecipeRepository := &mocks.RecipeRepository{}
			mockRecipeRepository.On("GetById", test.expectedId).
				Return(someRecipe, nil).
				Once()
			mockRecipeRepository.On("Update", someRecipe).
				Return(test.expectedErr).
				Once()

			underTest := &recipeInteractor{recipeRepository: mockRecipeRepository}

			actualErr := underTest.Update(test.expectedId, someRecipe)
			actualId := someRecipe.Id

			assert.Equal(t, test.expectedId, actualId)
			assert.Equal(t, test.expectedErr, actualErr)
			mockRecipeRepository.AssertExpectations(t)
		})
	}
}

func TestRecipeInteractor_Update_NotFound(t *testing.T) {
	someId := "someId"
	someRecipe := &entity.Recipe{Id: someId}

	expectedErr := errors.New("not found")

	mockRecipeRepository := &mocks.RecipeRepository{}
	mockRecipeRepository.On("GetById", someId).
		Return(nil, expectedErr).
		Once()

	underTest := &recipeInteractor{recipeRepository: mockRecipeRepository}

	actualErr := underTest.Update(someId, someRecipe)

	assert.Equal(t, expectedErr, actualErr)
	mockRecipeRepository.AssertExpectations(t)
	mockRecipeRepository.AssertNotCalled(t, "Update", someRecipe)
}

func TestRecipeInteractor_Delete(t *testing.T) {
	var tests = []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "deleted",
			expectedErr: nil,
		},
		{
			name:        "failed",
			expectedErr: errors.New("not deleted"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			someId := "someId"
			someRecipe := &entity.Recipe{Id: someId}

			mockRecipeRepository := &mocks.RecipeRepository{}
			mockRecipeRepository.On("GetById", someId).
				Return(someRecipe, nil).
				Once()
			mockRecipeRepository.On("Delete", someRecipe).
				Return(test.expectedErr).
				Once()

			underTest := &recipeInteractor{recipeRepository: mockRecipeRepository}

			actualErr := underTest.Delete(someId)

			assert.Equal(t, test.expectedErr, actualErr)
			mockRecipeRepository.AssertExpectations(t)
		})
	}
}

func TestRecipeInteractor_Delete_NotFound(t *testing.T) {
	someId := "someId"
	someRecipe := &entity.Recipe{Id: someId}

	expectedErr := errors.New("not found")

	mockRecipeRepository := &mocks.RecipeRepository{}
	mockRecipeRepository.On("GetById", someId).
		Return(nil, expectedErr).
		Once()

	underTest := &recipeInteractor{recipeRepository: mockRecipeRepository}

	actualErr := underTest.Delete(someId)

	assert.Equal(t, expectedErr, actualErr)
	mockRecipeRepository.AssertExpectations(t)
	mockRecipeRepository.AssertNotCalled(t, "Delete", someRecipe)
}

func TestRecipeInteractor_List(t *testing.T) {
	var tests = []struct {
		name            string
		expectedRecipes *[]entity.Recipe
		expectedErr     error
	}{
		{"found", &[]entity.Recipe{}, nil},
		{"not found", nil, errors.New("some error")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			somePage := 1918
			someLimit := 1992
			someFilter := "someFilter"

			mockRecipeRepository := &mocks.RecipeRepository{}
			mockRecipeRepository.On("GetAll", somePage, someLimit, someFilter).
				Return(test.expectedRecipes, test.expectedErr).
				Once()

			underTest := &recipeInteractor{recipeRepository: mockRecipeRepository}

			actualRecipes, actualErr := underTest.List(somePage, someLimit, someFilter)

			assert.Equal(t, test.expectedRecipes, actualRecipes)
			assert.Equal(t, test.expectedErr, actualErr)
			mockRecipeRepository.AssertExpectations(t)
		})
	}
}

func TestRecipeInteractor_Rate(t *testing.T) {
	var tests = []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "saved",
			expectedErr: nil,
		},
		{
			name:        "failed",
			expectedErr: errors.New("not saved"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedId := "someId"

			someRecipeId := "someRecipeId"
			someRecipe := &entity.Recipe{Id: someRecipeId}
			someRating := &entity.Rating{RecipeId: someRecipeId}

			mockIdGenerator := &mocks.IdGenerator{}
			mockIdGenerator.On("Generate").Return(expectedId).Once()

			mockRecipeRepository := &mocks.RecipeRepository{}
			mockRecipeRepository.On("GetById", someRecipeId).Return(someRecipe, nil).Once()

			mockRatingRepository := &mocks.RatingRepository{}
			mockRatingRepository.On("Save", someRating).Return(test.expectedErr).Once()

			underTest := &recipeInteractor{
				recipeRepository: mockRecipeRepository,
				ratingRepository: mockRatingRepository,
				idGenerator:      mockIdGenerator,
			}

			actualErr := underTest.Rate(someRecipeId, someRating)
			actualId := someRating.Id

			assert.Equal(t, expectedId, actualId)
			assert.Equal(t, test.expectedErr, actualErr)
			mockIdGenerator.AssertExpectations(t)
			mockRecipeRepository.AssertExpectations(t)
			mockRatingRepository.AssertExpectations(t)
		})
	}
}

func TestRecipeInteractor_Rate_NotFound(t *testing.T) {
	someId := "someId"
	someRating := &entity.Rating{}

	expectedErr := errors.New("not found")

	mockRecipeRepository := &mocks.RecipeRepository{}
	mockRecipeRepository.On("GetById", someId).
		Return(nil, expectedErr).
		Once()

	mockRatingRepository := &mocks.RatingRepository{}

	underTest := &recipeInteractor{recipeRepository: mockRecipeRepository}

	actualErr := underTest.Rate(someId, someRating)

	assert.Equal(t, expectedErr, actualErr)
	mockRecipeRepository.AssertExpectations(t)
	mockRatingRepository.AssertNotCalled(t, "Save", someRating)
}
