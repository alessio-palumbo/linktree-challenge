package links

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alessio-palumbo/linktree-challenge/handlers"
	"github.com/alessio-palumbo/linktree-challenge/middleware"
)

var (
	user1ID = "fac90185-d243-46f5-8797-e57ac9c2c293"
	user2ID = "9bce575b-1507-4a0f-a523-4072a72fc968"
	user3ID = "8c4664f7-ef96-48a1-80e3-804bbe0af06a"
	user4ID = "5bb13e12-3f40-42f0-be31-4a7ca007b432"
)

func TestIndexHandler_ServeHTTP(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	populateMockDB(mock)

	var testCases = []struct {
		name       string
		userID     string
		sortBy     string
		wantStatus int
		wantBody   *string
	}{
		{
			name:       "User with only classic links",
			userID:     user1ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "User with no links",
			userID:     user2ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "User with all types of links",
			userID:     user3ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Ordered request",
			userID:     user4ID,
			sortBy:     "created_at:asc",
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "https://linktree.com/api/links"
			if tc.sortBy != "" {
				url += fmt.Sprintf("?sort_by=%s", tc.sortBy)
			}

			req := httptest.NewRequest("GET", url, nil)
			req = middleware.CtxSetUserID(req.Context(), req, tc.userID)

			recorder := httptest.NewRecorder()

			IndexHandler(handlers.Group{DB: db}).ServeHTTP(recorder, req)

			if got := recorder.Code; got != tc.wantStatus {
				t.Errorf("got status %d, want %d", got, tc.wantStatus)
			}

			if tc.wantBody != nil {
				if got := recorder.Body.String(); got != *tc.wantBody {
					t.Errorf("got body %s, want %s", got, *tc.wantBody)
				}
			}
		})
	}
}

func populateMockDB(mock sqlmock.Sqlmock) {

	fields := []string{
		"l.id",
		"l.type",
		"l.title",
		"l.url",
		"l.thumbnail",
		"l.created_at",
		"sl.id",
		"sl.metadata",
	}

	// Set user1 mock DB. Only classic links
	user1Data := [][]driver.Value{
		[]driver.Value{
			"6e3060f3-4c99-41c7-a97b-a287399f3dd1",
			"classic",
			"First Link",
			"http://firstlink.com/1",
			nil,
			time.Now().UTC().Add(-24 * time.Hour),
			nil,
			nil,
		},
		[]driver.Value{
			"8d7a85a1-a875-49ad-9582-b8440e203650",
			"classic",
			"Second Link",
			"http://secondlink.com/2",
			nil,
			time.Now().UTC().Add(-8 * time.Hour),
			nil,
			nil,
		},
	}

	user1Rows := sqlmock.NewRows(fields)
	for _, row := range user1Data {
		user1Rows.AddRow(row...)
	}

	mock.ExpectQuery("SELECT l.id").WithArgs(user1ID).WillReturnRows(user1Rows)

	// Set user2 mock DB. No data
	mock.ExpectQuery("SELECT l.id").WithArgs(user2ID).WillReturnRows(sqlmock.NewRows(fields))

	// Set user3 mock DB. All types of links
	user3Data := [][]driver.Value{
		[]driver.Value{
			"e38c00bf-2187-4a75-ac4e-557cdfd8e263",
			"classic",
			"My Classic Link",
			"http://myclassiclink.com/classic",
			nil,
			time.Now().UTC().Add(-4 * time.Hour),
			nil,
			nil,
		},
		[]driver.Value{
			"f7265bc0-5d2f-43e3-b187-703239f798d4",
			"shows",
			"My Shows Link",
			nil,
			nil,
			time.Now().UTC().Add(-8 * time.Hour),
			"04e3c439-be86-4f19-ae1e-3f2bce732a41",
			[]byte(`{"id":"0ba388db-0a52-4979-97a2-f3c648e355e3","date":"Apr 01 2019",
			"venue":"Princess Theatre","location":"Melbourne","status":"sold-out"}`),
		},
		[]driver.Value{
			"f7265bc0-5d2f-43e3-b187-703239f798d4",
			"shows",
			"My Shows Link",
			nil,
			nil,
			time.Now().UTC().Add(-8 * time.Hour),
			"fb4ea9a5-8446-4201-a20b-818c944e3e09",
			[]byte(`{"id":"bff093b1-1857-4b74-94f1-d75fe8b44d41","date":"Sep 03 2020",
			"venue":"Opera House","location":"Sydney","status":"on-sale"}`),
		},
		[]driver.Value{
			"b626168a-6c34-44cb-bf94-667c76235a26",
			"music",
			"Music Link",
			"http://music-link.com/all-of-me",
			nil,
			time.Now().UTC().Add(-2 * time.Hour),
			"fbd19ca9-8006-448f-a2f0-52817ad7e9e1",
			[]byte(`{"name":"Spotify","url":"https://open.spotify.com/album/1YdXQgntClL3BhIXB0xpgs"}`),
		},
		[]driver.Value{
			"7fa60214-0827-45b6-b2f7-1690471760ad",
			"music",
			"Music Link",
			"http://music-link.com/all-of-me",
			nil,
			time.Now().UTC().Add(-2 * time.Hour),
			"2cbc2043-d67e-45fc-a687-7e147def358f",
			[]byte(`{"name":"SoundCloud","url":"https://soundcloud.com/johnlegend/all-of-me-3"}`),
		},
	}

	user3Rows := sqlmock.NewRows(fields)
	for _, row := range user3Data {
		user3Rows.AddRow(row...)
	}

	mock.ExpectQuery("SELECT l.id").WithArgs(user3ID).WillReturnRows(user3Rows)

	// Set user4 mock DB. Only classic links
	user4Data := [][]driver.Value{
		[]driver.Value{
			"a238c29f-d033-48a8-a20f-49c2756d25c4",
			"classic",
			"First Link",
			"http://firstlink.com/1",
			nil,
			time.Now().UTC().Add(-24 * time.Hour),
			nil,
			nil,
		},
		[]driver.Value{
			"a3c32ec1-54e0-4df4-b5f0-0be37b0a5b57",
			"classic",
			"Second Link",
			"http://secondlink.com/2",
			nil,
			time.Now().UTC().Add(-21 * time.Hour),
			nil,
			nil,
		},
	}

	user4Rows := sqlmock.NewRows(fields)
	for _, row := range user4Data {
		user1Rows.AddRow(row...)
	}

	mock.ExpectQuery("SELECT l.id").WithArgs(user4ID).WillReturnRows(user4Rows)
}

func Test_sortByClause(t *testing.T) {

	tests := []struct {
		name   string
		sortBy string
		want   string
	}{
		{
			name: "Empty string",
		},
		{
			name:   "Single column, no order",
			sortBy: "created_at",
			want:   "created_at asc",
		},
		{
			name:   "Single column, with order",
			sortBy: "created_at:asc",
			want:   "created_at asc",
		},
		{
			name:   "Single unknown column",
			sortBy: "order_id",
			want:   "",
		},
		{
			name:   "Single column, unknown order",
			sortBy: "created_at:descending",
			want:   "created_at asc",
		},
		{
			name:   "Multiple column, with order",
			sortBy: "created_at:desc,title:desc",
			want:   "created_at desc, title desc",
		},
		{
			name:   "Multiple column, no order",
			sortBy: "created_at,type",
			want:   "created_at asc, type asc",
		},
		{
			name:   "Invalid multiple column, with order",
			sortBy: "created_at:descending,order:asc",
			want:   "created_at asc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortByClause(tt.sortBy); got != tt.want {
				t.Errorf("sortByClause() = %v, want %v", got, tt.want)
			}
		})
	}
}
