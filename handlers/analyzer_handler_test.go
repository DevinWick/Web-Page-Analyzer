package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	router := gin.Default()
	router.LoadHTMLGlob("../pages/*")
	router.GET("/", IndexPathHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Web Page Analyzer")
}

func TestAnalyzeHandler(t *testing.T) {
	// Setup test server with mock responses
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/valid" {
			w.Write([]byte(`
				<!DOCTYPE html>
				<html>
					<head><title>Test Page</title></head>
					<body>
						<h1>Heading</h1>
						<a href="/internal">Internal</a>
						<a href="http://external.com">External</a>
					</body>
				</html>
			`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	router := gin.Default()
	router.LoadHTMLGlob("../pages/*")
	router.POST("/analyze", AnalyzeHandler)

	tests := []struct {
		name       string
		url        string
		statusCode int
		contains   string
	}{
		{
			name:       "Valid URL",
			url:        ts.URL + "/valid",
			statusCode: http.StatusOK,
			contains:   "Test Page",
		},
		{
			name:       "Invalid URL",
			url:        "invalid-url",
			statusCode: http.StatusBadRequest,
			contains:   "Invalid URL",
		},
		{
			name:       "Whitespace",
			url:        " ",
			statusCode: http.StatusBadRequest,
			contains:   "URL is required",
		},
		{
			name:       "Non-existent URL",
			url:        ts.URL + "/nonexistent",
			statusCode: http.StatusOK,
			contains:   "received non-200 status code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("url", tt.url)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/analyze", strings.NewReader(form.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.contains)
		})
	}
}
