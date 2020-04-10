package middleware

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	e "github.com/alessio-palumbo/linktree-challenge/errors"
)

var (
	token        = "15fe817c-72d7-49c1-bffc-0257dbd263e3"
	invalidToken = "05de14de-537a-4819-814a-85ec6c66dd35"
	userID       = "fac90185-d243-46f5-8797-e57ac9c2c293"

	expiredTimestamp = time.Now().UTC().Add(-1 * time.Minute)
	validTimestamp   = time.Now().UTC().Add(24 * time.Hour)
)

func TestAuth_ServeHTTP(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var testCases = []struct {
		name       string
		headers    map[string]string
		queryArgs  []driver.Value
		wantStatus int
		wantErr    string
		reqUID     string
	}{
		{
			name:       "Missing Authorization header",
			wantStatus: http.StatusUnauthorized,
			wantErr:    e.JSONError(errTokenMissing),
		},
		{
			name:       "Missing token in Authorization",
			headers:    map[string]string{"Authorization": "Bearer "},
			wantStatus: http.StatusUnauthorized,
			wantErr:    e.JSONError(errTokenMissing),
		},
		{
			name:       "Token not found",
			headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", invalidToken)},
			queryArgs:  []driver.Value{token, userID, expiredTimestamp},
			wantStatus: http.StatusUnauthorized,
			wantErr:    e.JSONError(errTokenInvalid),
		},
		{
			name:       "Token expired",
			headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
			queryArgs:  []driver.Value{token, userID, expiredTimestamp},
			wantStatus: http.StatusUnauthorized,
			wantErr:    e.JSONError(errTokenInvalid),
		},
		{
			name:       "Token valid",
			headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
			queryArgs:  []driver.Value{token, userID, validTimestamp},
			wantStatus: http.StatusOK,
			reqUID:     userID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			populate(mock, tc.queryArgs)

			url := url.URL{Scheme: "https", Host: "example.com", Path: "/api/links"}
			req := httptest.NewRequest("GET", url.String(), nil)

			for k, v := range tc.headers {
				req.Header.Add(k, v)
			}
			recorder := httptest.NewRecorder()

			var requestUserID string
			NewAuth(db).ServeHTTP(recorder, req, func(w http.ResponseWriter, r *http.Request) {
				requestUserID = CtxReqUserID(r.Context())
			})

			if got := recorder.Code; got != tc.wantStatus {
				t.Errorf("got status %d, want %d", got, tc.wantStatus)
			}

			if got := recorder.Body.String(); got != tc.wantErr {
				t.Errorf("got error %s, want %s", got, tc.wantErr)
			}

			if got := requestUserID; got != tc.reqUID {
				t.Errorf("got requesterID '%s', want '%s'", got, tc.reqUID)
			}

		})
	}
}

func populate(mock sqlmock.Sqlmock, qArgs []driver.Value) {
	q := "SELECT user_id"
	if len(qArgs) != 3 {
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		return
	}

	rows := sqlmock.NewRows([]string{"user_id"}).AddRow(qArgs[1])
	mock.ExpectQuery(q).WithArgs(qArgs[0]).WillReturnRows(rows)
}
