package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/handlers"
)

// TestRegisterRoutes tests the RegisterRoutes function.
func TestRegisterRoutes(t *testing.T) {
	r := chi.NewRouter()

	// Initialize a mock or simple instance of RequestHandler
	requestHandler := &handlers.RequestHandler{}

	RegisterRoutes(r, requestHandler)

	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"GET", "/api/course", http.StatusOK},
		{"GET", "/api/course/1", http.StatusOK},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.expectedCode, rr.Code)
	}
}
