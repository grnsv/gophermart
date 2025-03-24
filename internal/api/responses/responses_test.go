package responses

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestWriteJSON_Success(t *testing.T) {
	recorder := httptest.NewRecorder()
	response := Response{Data: "test data"}

	WriteJSON(recorder, http.StatusOK, response)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var result Response
	err := json.NewDecoder(recorder.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, response, result)
}

func TestWriteJSON_Error(t *testing.T) {
	recorder := httptest.NewRecorder()
	response := ErrorResponse{Message: "test error"}

	WriteJSON(recorder, http.StatusBadRequest, response)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var result ErrorResponse
	err := json.NewDecoder(recorder.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, response, result)
}

func TestNewErrorsResponse(t *testing.T) {
	validate := validator.New()
	type TestStruct struct {
		Field1 string `validate:"required"`
		Field2 int    `validate:"min=10"`
	}

	testData := TestStruct{}
	err := validate.Struct(testData)

	errorsResponse := NewErrorsResponse("Validation failed", err)

	assert.Equal(t, "Validation failed", errorsResponse.Message)
	assert.Contains(t, errorsResponse.Errors, "Field1")
	assert.Equal(t, "required", errorsResponse.Errors["Field1"])
	assert.Contains(t, errorsResponse.Errors, "Field2")
	assert.Equal(t, "min", errorsResponse.Errors["Field2"])
}
