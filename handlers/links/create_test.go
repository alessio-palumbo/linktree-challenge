package links

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alessio-palumbo/linktree-challenge/handlers"
	"github.com/alessio-palumbo/linktree-challenge/middleware"
	"github.com/alessio-palumbo/linktree-challenge/validator"
)

func TestPostHandler_ServeHTTP(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var testCases = []struct {
		name       string
		userID     string
		payload    string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Invalid payload, missing type",
			userID:     user1ID,
			payload:    `{"title":"first link"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"validation errors: Type is required"}`,
		},
		{
			name:       "Invalid payload, title is over 144 characters",
			userID:     user1ID,
			payload:    fmt.Sprintf(`{"type":"classic","title":"%s"}`, strings.Repeat("a", 145)),
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"validation errors: Title is longer than 144 characters"}`,
		},
		{
			name:       "Invalid type",
			userID:     user1ID,
			payload:    fmt.Sprintf(`{"type":"classic","title":"%s"}`, strings.Repeat("a", 145)),
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"validation errors: Title is longer than 144 characters"}`,
		},
		{
			name:       "Music link with sublinks",
			userID:     user1ID,
			payload:    `{"type":"music","sublinks":[{"name":"Spotify","url":"http://music-link.com/all-of-me"}]}`,
			wantStatus: http.StatusOK,
			wantBody: `{"id":"","type":"music","title":null,"url":null,"sublinks":[{"id":"",` +
				`"name":"Spotify","url":"http://music-link.com/all-of-me"}]}`,
		},
		{
			name:   "Music link with show sublinks but with valid fields",
			userID: user1ID,
			payload: `{"type":"music","sublinks":[{"date":"Apr 01 2019","name":"Cats",` +
				`"venue":"Princess Theatre","location":"Melbourne","status": "sold-out",` +
				`"url":"https://cats.com.au"}]}`,
			wantStatus: http.StatusOK,
			wantBody: `{"id":"","type":"music","title":null,"url":null,"sublinks":[{` +
				`"id":"","name":"Cats","url":"https://cats.com.au"}]}`,
		},
		{
			name:       "Music link with missing required fields",
			userID:     user1ID,
			payload:    `{"type":"music","sublinks":[{"name":"Spotify"}]}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"validation errors: URL is required"}`,
		},
		{
			name:   "Show link with valid sublink fields",
			userID: user1ID,
			payload: `{"type":"shows","sublinks":[{"date":"Apr 01 2019","name":"Cats",` +
				`"venue":"Princess Theatre","location":"Melbourne","status": "sold-out",` +
				`"url":"https://cats.com.au"}]}`,
			wantStatus: http.StatusOK,
			wantBody: `{"id":"","type":"shows","title":null,"url":null,"sublinks":[{` +
				`"id":"","date":"Apr 01 2019","name":"Cats","venue":"Princess Theatre",` +
				`"location":"Melbourne","status":"sold-out","url":"https://cats.com.au"}]}`,
		},
		{
			name:   "Show link with invalid sublink fields",
			userID: user1ID,
			payload: `{"type":"shows","sublinks":[{"date":"Apr 31 2019","name":"Cats",` +
				`"status": "coming-soon","url":"https://cats.com.au"}]}`,
			wantStatus: http.StatusBadRequest,
			wantBody: `{"error":"validation errors: Date is invalid, Venue is required ` +
				`in absence of Location, Location is required in absence of Venue, Status is invalid"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			req := httptest.NewRequest("POST", "https://linktree.com/api/links", strings.NewReader(tc.payload))
			req = middleware.CtxSetUserID(req.Context(), req, tc.userID)

			recorder := httptest.NewRecorder()

			PostHandler(handlers.Group{DB: db, Validator: validator.New()}).ServeHTTP(recorder, req)

			if got := recorder.Code; got != tc.wantStatus {
				t.Errorf("got status %d, want %d", got, tc.wantStatus)
			}

			if tc.wantBody != "" {
				if got := strings.TrimSpace(recorder.Body.String()); got != tc.wantBody {
					t.Errorf("got body %s, want %s", got, tc.wantBody)
				}
			}
		})
	}
}
