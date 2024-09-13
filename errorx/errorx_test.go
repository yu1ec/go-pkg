package errorx_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/yu1ec/go-pkg/errorx"
)

func TestNewError(t *testing.T) {
	err := errorx.NewError(errorx.ErrBadRequest, "ERR001", "Invalid input")
	if err.HttpStatusCode != errorx.ErrBadRequest {
		t.Errorf("Expected HTTP status code %d, got %d", errorx.ErrBadRequest, err.HttpStatusCode)
	}
	if err.ErrorCode != "ERR001" {
		t.Errorf("Expected error code ERR001, got %s", err.ErrorCode)
	}
	if err.Reason != "Invalid input" {
		t.Errorf("Expected reason 'Invalid input', got '%s'", err.Reason)
	}
}

func TestError_Error(t *testing.T) {
	err := errorx.NewError(errorx.ErrBadRequest, "ERR001", "Invalid input")
	if err.Error() != "Invalid input" {
		t.Errorf("Expected error message 'Invalid input', got '%s'", err.Error())
	}
}

func TestError_WithCause(t *testing.T) {
	err := errorx.NewError(errorx.ErrBadRequest, "ERR001", "Invalid input")
	newErr := err.WithCause("Missing field")
	if newErr.Error() != "Invalid input: Missing field" {
		t.Errorf("Expected error message 'Invalid input: Missing field', got '%s'", newErr.Error())
	}
}

func TestWithCause(t *testing.T) {
	t.Run("With errorx.Error", func(t *testing.T) {
		err := errorx.NewError(errorx.ErrBadRequest, "ERR001", "Invalid input")
		newErr := errorx.WithCause(err, "Missing field")
		if newErr.Error() != "Invalid input: Missing field" {
			t.Errorf("Expected error message 'Invalid input: Missing field', got '%s'", newErr.Error())
		}
	})

	t.Run("With standard error", func(t *testing.T) {
		err := errors.New("Standard error")
		newErr := errorx.WithCause(err, "Additional info")
		if newErr.Error() != "Additional info: Standard error" {
			t.Errorf("Expected error message 'Additional info: Standard error', got '%s'", newErr.Error())
		}
	})
}

func TestError_Data(t *testing.T) {
	err := errorx.NewError(errorx.ErrBadRequest, "ERR001", "Invalid input")
	statusCode, respErr := err.Data()

	if statusCode != errorx.ErrBadRequest {
		t.Errorf("Expected HTTP status code %d, got %d", errorx.ErrBadRequest, statusCode)
	}

	if respErr.Code != "ERR001" {
		t.Errorf("Expected error code ERR001, got %s", respErr.Code)
	}

	if respErr.Reason != "Invalid input" {
		t.Errorf("Expected reason 'Invalid input', got '%s'", respErr.Reason)
	}
}

func TestHttpStatusCodes(t *testing.T) {
	testCases := []struct {
		code     errorx.HttpStatusCode
		expected int
	}{
		{errorx.ErrBadRequest, 400},
		{errorx.ErrUnauthorized, 401},
		{errorx.ErrForbidden, 403},
		{errorx.ErrNotFound, 404},
		{errorx.ErrMethodNotAllowed, 405},
		{errorx.ErrNotAcceptable, 406},
		{errorx.ErrInternalServerError, 500},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("HttpStatusCode_%d", tc.expected), func(t *testing.T) {
			if int(tc.code) != tc.expected {
				t.Errorf("Expected HTTP status code %d, got %d", tc.expected, tc.code)
			}
		})
	}
}
