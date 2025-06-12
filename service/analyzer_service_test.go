package service

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/devinwick/web-page-analyzer/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"net/url"
	"strings"
	"testing"
)

func TestDetermineHTMLVersion(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "HTML5",
			html:     `<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`,
			expected: "HTML5",
		},
		{
			name:     "HTML 4.01 Strict",
			html:     `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd"><html><head>    <title>HTML 4 Sample Page</title></head><body></body></html>`,
			expected: "HTML 4.01",
		},
		{
			name:     "XHTML 1.0",
			html:     `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"><html xmlns="http://www.w3.org/1999/xhtml"><head><title>Test</title></head><body></body></html>`,
			expected: "XHTML 1.0",
		},
		{
			name:     "No DOCTYPE",
			html:     `<html><head><title>Test</title></head><body></body></html>`,
			expected: "HTML5 (no DOCTYPE found, assuming HTML5)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootNode, err := html.Parse(strings.NewReader((tt.html)))
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, determineHTMLVersion(rootNode))
		})
	}
}

func TestAnalyzeHeadings(t *testing.T) {
	htmlStr := `
	<html>
		<head><title>Test</title></head>
		<body>
			<h1>Heading 1</h1>
			<h2>Heading 2</h2>
			<h2>Heading 2</h2>
			<h3>Heading 3</h3>
		</body>
	</html>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	assert.NoError(t, err)

	result := &model.AnalysisResult{Headings: make(map[string]int)}
	analyzeHeadings(doc, result)

	assert.Equal(t, 1, result.Headings["h1"])
	assert.Equal(t, 2, result.Headings["h2"])
	assert.Equal(t, 1, result.Headings["h3"])
	assert.Equal(t, 0, result.Headings["h4"])
}

func TestCheckForLoginForm(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name:     "Password input present",
			html:     `<form><input type="password" name="pass"></form>`,
			expected: true,
		},
		{
			name:     "Login in form ID",
			html:     `<form id="loginForm"></form>`,
			expected: true,
		},
		{
			name:     "Login in form class",
			html:     `<form class="user-login"></form>`,
			expected: true,
		},
		{
			name:     "Login form with username",
			html:     `<form><input type="email" name="username"></form>`,
			expected: true,
		},
		{
			name:     "No login form",
			html:     `<form><input type="text" name="firstName"></form>`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, checkForLoginForm(doc))
		})
	}
}

func TestIsInternalLink(t *testing.T) {
	baseURL, _ := url.Parse("http://google.com")

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{"Same domain", "http://google.com/page", true},
		{"Relative path", "/page", true},
		{"Different domain", "http://test.com", false},
		{"Subdomain", "http://info.google.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, _ := url.Parse(tt.link)
			assert.Equal(t, tt.expected, isInternalLink(baseURL, target))
		})
	}
}
