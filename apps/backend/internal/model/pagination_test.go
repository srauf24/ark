package model

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

// Test 1: TestPaginationParams_SetDefaults_ZeroLimit
func TestPaginationParams_SetDefaults_ZeroLimit(t *testing.T) {
	params := PaginationParams{Limit: 0, Offset: 10}
	params.SetDefaults(20)

	if params.Limit != 20 {
		t.Errorf("Expected Limit to be 20, got %d", params.Limit)
	}
	if params.Offset != 10 {
		t.Errorf("Expected Offset to stay 10, got %d", params.Offset)
	}
}

// Test 2: TestPaginationParams_SetDefaults_NegativeOffset
func TestPaginationParams_SetDefaults_NegativeOffset(t *testing.T) {
	params := PaginationParams{Limit: 10, Offset: -5}
	params.SetDefaults(20)

	if params.Limit != 10 {
		t.Errorf("Expected Limit to stay 10, got %d", params.Limit)
	}
	if params.Offset != 0 {
		t.Errorf("Expected Offset to become 0, got %d", params.Offset)
	}
}

// Test 3: TestPaginationParams_SetDefaults_BothZeroAndNegative
func TestPaginationParams_SetDefaults_BothZeroAndNegative(t *testing.T) {
	params := PaginationParams{Limit: 0, Offset: -10}
	params.SetDefaults(50)

	if params.Limit != 50 {
		t.Errorf("Expected Limit to become 50, got %d", params.Limit)
	}
	if params.Offset != 0 {
		t.Errorf("Expected Offset to become 0, got %d", params.Offset)
	}
}

// Test 4: TestPaginationParams_SetDefaults_ValidValues
func TestPaginationParams_SetDefaults_ValidValues(t *testing.T) {
	params := PaginationParams{Limit: 30, Offset: 20}
	params.SetDefaults(50)

	if params.Limit != 30 {
		t.Errorf("Expected Limit to stay 30, got %d", params.Limit)
	}
	if params.Offset != 20 {
		t.Errorf("Expected Offset to stay 20, got %d", params.Offset)
	}
}

// Test 5: TestPaginationParams_Validation_ValidParams
func TestPaginationParams_Validation_ValidParams(t *testing.T) {
	validate := validator.New()
	params := PaginationParams{Limit: 50, Offset: 0}

	err := validate.Struct(params)
	if err != nil {
		t.Errorf("Expected validation to pass, got error: %v", err)
	}
}

// Test 6: TestPaginationParams_Validation_LimitTooHigh
func TestPaginationParams_Validation_LimitTooHigh(t *testing.T) {
	validate := validator.New()
	params := PaginationParams{Limit: 101, Offset: 0}

	err := validate.Struct(params)
	if err == nil {
		t.Error("Expected validation error for Limit > 100, got none")
	}

	// Verify the error is on the Limit field
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Limit" && fieldError.Tag() == "max" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Limit field with 'max' tag")
		}
	}
}

// Test 7: TestPaginationParams_Validation_LimitNegative
func TestPaginationParams_Validation_LimitNegative(t *testing.T) {
	validate := validator.New()
	params := PaginationParams{Limit: -1, Offset: 0}

	err := validate.Struct(params)
	if err == nil {
		t.Error("Expected validation error for Limit < 0, got none")
	}

	// Verify the error is on the Limit field
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Limit" && fieldError.Tag() == "min" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Limit field with 'min' tag")
		}
	}
}

// Test 8: TestPaginationParams_Validation_OffsetNegative
func TestPaginationParams_Validation_OffsetNegative(t *testing.T) {
	validate := validator.New()
	params := PaginationParams{Limit: 20, Offset: -1}

	err := validate.Struct(params)
	if err == nil {
		t.Error("Expected validation error for Offset < 0, got none")
	}

	// Verify the error is on the Offset field
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Offset" && fieldError.Tag() == "min" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Offset field with 'min' tag")
		}
	}
}

// Test 9: TestPaginationParams_Validation_ZeroValuesValid
func TestPaginationParams_Validation_ZeroValuesValid(t *testing.T) {
	validate := validator.New()
	params := PaginationParams{Limit: 0, Offset: 0}

	err := validate.Struct(params)
	if err != nil {
		t.Errorf("Expected validation to pass for zero values (omitempty), got error: %v", err)
	}
}
