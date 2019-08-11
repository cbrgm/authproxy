package errors

import (
	"testing"
	"net/http"
	"errors"
)

func TestErrorNew(t *testing.T) {
	if !IsInternalError(NewInternalError(errors.New("message"))) {
		t.Errorf("expected to be %v", http.StatusInternalServerError)
	}

	if !IsUnauthorized(NewUnauthorized("message")) {
		t.Errorf("expected to be %v", http.StatusUnauthorized)
	}
}
