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

// PaginationMeta Tests

// Test 10: TestNewPaginationMeta_FirstPage
func TestNewPaginationMeta_FirstPage(t *testing.T) {
	meta := NewPaginationMeta(100, 20, 0)

	if meta.Total != 100 {
		t.Errorf("Expected Total to be 100, got %d", meta.Total)
	}
	if meta.Limit != 20 {
		t.Errorf("Expected Limit to be 20, got %d", meta.Limit)
	}
	if meta.Offset != 0 {
		t.Errorf("Expected Offset to be 0, got %d", meta.Offset)
	}
	if !meta.HasNext {
		t.Error("Expected HasNext to be true on first page with more results")
	}
	if meta.HasPrev {
		t.Error("Expected HasPrev to be false on first page")
	}
}

// Test 11: TestNewPaginationMeta_MiddlePage
func TestNewPaginationMeta_MiddlePage(t *testing.T) {
	meta := NewPaginationMeta(100, 20, 40)

	if meta.Total != 100 {
		t.Errorf("Expected Total to be 100, got %d", meta.Total)
	}
	if !meta.HasNext {
		t.Error("Expected HasNext to be true on middle page")
	}
	if !meta.HasPrev {
		t.Error("Expected HasPrev to be true on middle page")
	}
}

// Test 12: TestNewPaginationMeta_LastPage
func TestNewPaginationMeta_LastPage(t *testing.T) {
	meta := NewPaginationMeta(100, 20, 80)

	if meta.HasNext {
		t.Error("Expected HasNext to be false on last page")
	}
	if !meta.HasPrev {
		t.Error("Expected HasPrev to be true on last page")
	}
}

// Test 13: TestNewPaginationMeta_ExactPageBoundary
func TestNewPaginationMeta_ExactPageBoundary(t *testing.T) {
	meta := NewPaginationMeta(100, 20, 60)

	// offset 60 + limit 20 = 80, which is still < 100
	if !meta.HasNext {
		t.Error("Expected HasNext to be true when offset+limit < total")
	}
	if !meta.HasPrev {
		t.Error("Expected HasPrev to be true when offset > 0")
	}
}

// Test 14: TestNewPaginationMeta_SinglePage
func TestNewPaginationMeta_SinglePage(t *testing.T) {
	meta := NewPaginationMeta(15, 20, 0)

	if meta.HasNext {
		t.Error("Expected HasNext to be false when all results fit on one page")
	}
	if meta.HasPrev {
		t.Error("Expected HasPrev to be false on first page")
	}
}

// Test 15: TestNewPaginationMeta_EmptyResults
func TestNewPaginationMeta_EmptyResults(t *testing.T) {
	meta := NewPaginationMeta(0, 20, 0)

	if meta.Total != 0 {
		t.Errorf("Expected Total to be 0, got %d", meta.Total)
	}
	if meta.HasNext {
		t.Error("Expected HasNext to be false with no results")
	}
	if meta.HasPrev {
		t.Error("Expected HasPrev to be false with no results")
	}
}

// Test 16: TestNewPaginationMeta_LargeOffset
func TestNewPaginationMeta_LargeOffset(t *testing.T) {
	meta := NewPaginationMeta(100, 20, 200)

	// offset 200 is beyond total 100
	if meta.HasNext {
		t.Error("Expected HasNext to be false when offset exceeds total")
	}
	if !meta.HasPrev {
		t.Error("Expected HasPrev to be true when offset > 0")
	}
}

// Test 17: TestNewPaginationMeta_FieldValues
func TestNewPaginationMeta_FieldValues(t *testing.T) {
	meta := NewPaginationMeta(100, 20, 40)

	if meta.Total != 100 {
		t.Errorf("Expected Total=100, got %d", meta.Total)
	}
	if meta.Limit != 20 {
		t.Errorf("Expected Limit=20, got %d", meta.Limit)
	}
	if meta.Offset != 40 {
		t.Errorf("Expected Offset=40, got %d", meta.Offset)
	}
}
