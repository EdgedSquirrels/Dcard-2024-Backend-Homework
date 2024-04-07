package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	name         string
	method       string
	url          string
	bodyPath     string
	expectedCode int
	expectedBody string
}

func testAPI(t *testing.T, router *gin.Engine, req *http.Request, test TestCase, wg *sync.WaitGroup) {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, test.expectedCode, w.Code)
	if test.expectedBody != "" {
		assert.Equal(t, test.expectedBody, w.Body.String())
	}
	wg.Done()
}

func TestPostGetAds(t *testing.T) {
	router := setupRouter()

	testCases := []TestCase{
		{"get empty ads", "GET", "/api/v1/ad", "", http.StatusOK, `{"items":[]}`},
		{"post AD 1", "POST", "/api/v1/ad", "test/AD 1.json", http.StatusOK, ""},
		{"post AD 2", "POST", "/api/v1/ad", "test/AD 2.json", http.StatusOK, ""},
		{"get ads", "GET", "/api/v1/ad?age=39", "", http.StatusOK, `{"items":[{"title":"AD 2","endAt":"2024-12-31T16:00:00Z"}]}`},
		{"post AD 3", "POST", "/api/v1/ad", "test/AD 3.json", http.StatusBadRequest, ""},
		{"post AD 4", "POST", "/api/v1/ad", "test/AD 4.json", http.StatusOK, ""},
		{"post AD 5", "POST", "/api/v1/ad", "test/AD 5.json", http.StatusOK, ""},
		{"post AD 6", "POST", "/api/v1/ad", "test/AD 6.json", http.StatusBadRequest, ""},
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
			req, err = http.NewRequest(testCase.method, testCase.url, f)
			checkErr(err)
		}
		router.ServeHTTP(w, req)
		assert.Equal(t, testCase.expectedCode, w.Code)
		if testCase.expectedBody != "" {
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		}
	}

	// send the last test 1000 times to check the performance
	var wg sync.WaitGroup
	test_num := 1000
	wg.Add(test_num)
	test := testCases[len(testCases)-1]
	req, _ := http.NewRequest(test.method, test.url, nil)
	for i := 0; i < test_num; i++ {
		go testAPI(t, router, req, test, &wg)
	}
	wg.Wait()
}
