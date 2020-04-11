package validator

import (
	"testing"

	validator "gopkg.in/go-playground/validator.v9"
)

func TestCustomValidator_Validate(t *testing.T) {

	testCases := []struct {
		name            string
		payload         interface{}
		wantErr         bool
		wantTranslation string
	}{
		{
			name: "Validate missing required fields",
			payload: struct {
				ID string `validate:"required"`
			}{
				ID: "",
			},
			wantErr:         true,
			wantTranslation: "validation errors: ID is required",
		},
		{
			name: "Validate correct required fields",
			payload: struct {
				ID string `validate:"required"`
			}{
				ID: "00001",
			},
			wantErr: false,
		},
		{
			name: "Missing required without field",
			payload: struct {
				ID   string `validate:"required_without=Name"`
				Name string
				Age  int
			}{
				Age: 43,
			},
			wantErr:         true,
			wantTranslation: "validation errors: ID is required in absence of Name",
		},
		{
			name: "Correct required without field",
			payload: struct {
				ID   string `validate:"required_without=Name"`
				Name string
				Age  int
			}{
				Name: "John",
				Age:  43,
			},
			wantErr: false,
		},
		{
			name: "Max number reached",
			payload: struct {
				Age int `validate:"max=42"`
			}{
				Age: 43,
			},
			wantErr:         true,
			wantTranslation: "validation errors: Age is greater than 42",
		},
		{
			name: "Max string length reached",
			payload: struct {
				Title string `validate:"max=10"`
			}{
				Title: "hello world",
			},
			wantErr:         true,
			wantTranslation: "validation errors: Title is longer than 10 characters",
		},
	}

	cv := &CustomValidator{
		validator: validator.New(),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cv.Validate(tc.payload)
			if (err != nil) != tc.wantErr {
				t.Errorf("CustomValidator.Validate() error = %v, wantErr %v", err, tc.wantErr)
			}

			if err != nil && tc.wantTranslation != "" && err.Error() != tc.wantTranslation {
				t.Errorf("CustomValidator.Validate() translation error, got %v, wantTranslation %v", err.Error(), tc.wantTranslation)
			}
		})
	}
}
