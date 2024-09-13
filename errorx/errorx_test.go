package errorx

import (
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(ErrBadRequest, "INVALID_INPUT", "Invalid input")

	if err.HttpStatusCode() != int(ErrBadRequest) {
		t.Errorf("Expected HTTP status code %d, got %d", int(ErrBadRequest), err.HttpStatusCode())
	}

	if err.ErrorCode() != "INVALID_INPUT" {
		t.Errorf("Expected error code 'INVALID_INPUT', got '%s'", err.ErrorCode())
	}

	if err.Error() != "Invalid input" {
		t.Errorf("Expected reason 'Invalid input', got '%s'", err.Error())
	}
}

func TestWithCause(t *testing.T) {
	originalErr := NewError(ErrBadRequest, "INVALID_INPUT", "Invalid input")
	err := WithCause(originalErr, "Additional information")

	if err.Error() != "Invalid input: Additional information" {
		t.Errorf("Expected 'Invalid input: Additional information', got '%s'", err.Error())
	}
}

func TestErrorData(t *testing.T) {
	err := NewError(ErrBadRequest, "INVALID_INPUT", "Invalid input")
	data := err.Data()

	if data.Code != "INVALID_INPUT" {
		t.Errorf("Expected code 'INVALID_INPUT', got '%s'", data.Code)
	}

	if data.Reason != "Invalid input" {
		t.Errorf("Expected reason 'Invalid input', got '%s'", data.Reason)
	}
}
