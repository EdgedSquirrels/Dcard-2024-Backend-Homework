package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func DeleteAds() {
	db, err := sql.Open("postgres", "dbname=ad sslmode=disable")
	checkErr(err)

	_, err = db.Exec("DELETE FROM ad")
	checkErr(err)
}

type TestCase struct {
	name         string
	method       string
	url          string
	bodyPath     string
	expectedCode int
	expectedBody string
}

func TestPostGetAds(t *testing.T) {
	DeleteAds()
	router := setupRouter()

	testCases := []TestCase{
		{"get empty ads", "GET", "/api/v1/ad", "", http.StatusOK, `{"items":[]}`},
		{"post AD 1", "POST", "/api/v1/ad", "test/AD 1.json", http.StatusOK, ""},
		{"post AD 2", "POST", "/api/v1/ad", "test/AD 2.json", http.StatusOK, ""},
		{"get ads", "GET", "/api/v1/ad?age=39", "", http.StatusOK, `{"items":[{"title":"AD 2","endAt":"2024-12-31T16:00:00Z"}]}`},
		{"post AD 3", "POST", "/api/v1/ad", "test/AD 3.json", http.StatusBadRequest, ""},
		{"post AD 4", "POST", "/api/v1/ad", "test/AD 4.json", http.StatusOK, ""},
		{"post AD 5", "POST", "/api/v1/ad", "test/AD 5.json", http.StatusOK, ""},
		{"post AD 6", "POST", "/api/v1/ad", "test/AD 3.json", http.StatusBadRequest, ""},
		{"get ads", "GET", "/api/v1/ad?age=30", "", http.StatusOK, `{"items":[{"title":"AD 5","endAt":"2024-12-31T15:00:00Z"},{"title":"AD 2","endAt":"2024-12-31T16:00:00Z"},{"title":"AD 3","endAt":"2024-12-31T16:00:00Z"}]}`},
		{"get ads", "GET", "/api/v1/ad?age=50", "", http.StatusOK, `{"items":[{"title":"AD 2","endAt":"2024-12-31T16:00:00Z"}]}`},
		{"get ads", "GET", "/api/v1/ad?platform=ios&limit=1&offset=1", "", http.StatusOK, `{"items":[{"title":"AD 2","endAt":"2024-12-31T16:00:00Z"}]}`},
		{"get ads", "GET", "/api/v1/ad?platform=idk&limit=1&offset=1", "", http.StatusBadRequest, ""},
		{"get ads", "GET", "/api/v1/ad?age=30&platform=android", "", http.StatusOK, `{"items":[{"title":"AD 2","endAt":"2024-12-31T16:00:00Z"},{"title":"AD 3","endAt":"2024-12-31T16:00:00Z"}]}`},
	}

	for _, testCase := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(testCase.method, testCase.url, nil)
		if testCase.bodyPath != "" {
			f, err := os.Open(testCase.bodyPath)
			checkErr(err)
			req, _ = http.NewRequest(testCase.method, testCase.url, f)
		}
		router.ServeHTTP(w, req)
		assert.Equal(t, testCase.expectedCode, w.Code)
		if testCase.expectedBody != "" {
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		}
	}
}
