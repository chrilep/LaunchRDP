package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI(t *testing.T) {
	server, err := NewServer(8080)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Test an API endpoint", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/users", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.handleUsers)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})
}
