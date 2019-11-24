package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SananGuliyev/goddd/domain/entity"
	"github.com/SananGuliyev/goddd/domain/throwable"
	"github.com/SananGuliyev/goddd/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecipeHandler_Get(t *testing.T) {
	someRecipe := &entity.Recipe{Id: "someExistingId"}
	someRating := float32(5)
	someMessage := "Some not found message"
	someErr := throwable.NewNotFound(someMessage)

	expectedResponse := struct {
		Recipe *entity.Recipe `json:"recipe"`
		Rating float32        `json:"rating"`
	}{Recipe: someRecipe, Rating: someRating}

	tests := []struct {
		name               string
		someId             string
		someRecipe         *entity.Recipe
		someRating         float32
		expectedResponse   interface{}
		expectedErr        error
		expectedStatusCode int
	}{
		{
			name:               "success case",
			someId:             "someExistingId",
			someRecipe:         someRecipe,
			someRating:         someRating,
			expectedResponse:   &SuccessResponse{StatusCode: http.StatusOK, Data: expectedResponse},
			expectedErr:        nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "not found case",
			someId:             "someNonExistingId",
			someRecipe:         nil,
			someRating:         float32(0),
			expectedResponse:   &ErrorResponse{StatusCode: http.StatusNotFound, Message: someMessage},
			expectedErr:        someErr,
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target := fmt.Sprintf("https://goddd.com/recipes/%s", test.someId)
			expectedResponse := test.expectedResponse

			expected, _ := json.Marshal(expectedResponse)

			request := httptest.NewRequest(http.MethodGet, target, nil)
			writer := httptest.NewRecorder()

			mockRecipeInteractor := &mocks.RecipeInteractor{}
			mockRecipeInteractor.On("Get", test.someId).
				Return(test.someRecipe, test.someRating, test.expectedErr).
				Once()

			underTest := RecipeHandler{Interactor: mockRecipeInteractor}
			underTest.Get(writer, request, test.someId)

			response := writer.Result()
			actual, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)
			assert.Equal(t, expected, actual)
			mockRecipeInteractor.AssertExpectations(t)
		})
	}
}

func TestRecipeHandler_Create(t *testing.T) {
	someRecipe := &entity.Recipe{}

	someMessage := "some database issue"
	someErr := errors.New(someMessage)

	tests := []struct {
		name               string
		someRecipe         *entity.Recipe
		expectedResponse   interface{}
		expectedErr        error
		expectedStatusCode int
	}{
		{
			name:               "success case",
			someRecipe:         someRecipe,
			expectedResponse:   &SuccessResponse{StatusCode: http.StatusCreated, Data: someRecipe},
			expectedErr:        nil,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "interactor failed case",
			someRecipe:         &entity.Recipe{},
			expectedResponse:   &ErrorResponse{StatusCode: http.StatusInternalServerError, Message: someMessage},
			expectedErr:        someErr,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedResponse := test.expectedResponse
			expected, _ := json.Marshal(expectedResponse)

			someRequestBody, _ := json.Marshal(test.someRecipe)

			request := httptest.NewRequest(
				http.MethodPost,
				"https://goddd.com/recipes",
				bytes.NewReader(someRequestBody),
			)
			writer := httptest.NewRecorder()

			mockRecipeInteractor := &mocks.RecipeInteractor{}
			mockRecipeInteractor.On("Create", test.someRecipe).Return(test.expectedErr).Once()

			underTest := RecipeHandler{Interactor: mockRecipeInteractor}
			underTest.Create(writer, request)

			response := writer.Result()
			actual, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)
			assert.Equal(t, expected, actual)
			mockRecipeInteractor.AssertExpectations(t)
		})
	}
}

func TestRecipeHandler_Create_WithInvalidBody(t *testing.T) {
	someInvalidRequestBody := []byte("{\"prepare_time\": \"someInvalidValue\"}")

	someMessage := "Request body is invalid"

	expectedResponse := &ErrorResponse{StatusCode: http.StatusBadRequest, Message: someMessage}
	expected, _ := json.Marshal(expectedResponse)

	request := httptest.NewRequest(
		http.MethodPost,
		"https://goddd.com/recipes",
		bytes.NewReader(someInvalidRequestBody),
	)
	writer := httptest.NewRecorder()

	mockRecipeInteractor := &mocks.RecipeInteractor{}

	underTest := RecipeHandler{Interactor: mockRecipeInteractor}
	underTest.Create(writer, request)

	response := writer.Result()
	actual, _ := ioutil.ReadAll(response.Body)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, expected, actual)
	mockRecipeInteractor.AssertExpectations(t)
	mockRecipeInteractor.AssertNotCalled(t, "Create", &entity.Recipe{})
}
