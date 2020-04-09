package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alessio-palumbo/linktree-challenge/handlers"
)

func TestNew(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ok"}).AddRow(1)
	mock.ExpectQuery("SELECT true as ok").WillReturnRows(rows)

	var testCases = []struct {
		name         string
		token        string
		method, path string
		want         int
	}{
		{"status internal", "__TOKEN__", "GET", "/healthcheck", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := url.URL{Scheme: "https", Host: "example.com", Path: tc.path}
			req := httptest.NewRequest(tc.method, url.String(), nil)
			req.Header.Set("Authorization", "Bearer "+tc.token)
			recorder := httptest.NewRecorder()

			New(handlers.Group{DB: db}).ServeHTTP(recorder, req)

			if got := recorder.Code; got != tc.want {
				t.Errorf("got status %d, want %d", got, tc.want)
			}
		})
	}
}
