package service

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	LOGGER "github.com/devinwick/web-page-analyzer/logger"
	"github.com/devinwick/web-page-analyzer/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const webPageURLTimeout = 10 * time.Second
const linkTimeout = 10 * time.Second

var logger *logrus.Entry = LOGGER.Log.WithField("pkg", "service")

func AnalyzeWebPage(targetURL string) (*model.AnalysisResult, error) {
	result := &model.AnalysisResult{
		URL:      targetURL,
		Headings: make(map[string]int),
	}

	// get the web page
	client := &http.Client{Timeout: webPageURLTimeout}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		result.StatusCode = 500
		return result, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		result.StatusCode = 500
		return result, err
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("web page returned status: %d", resp.StatusCode)
	}

	rootNode, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	trackTime("determineHTMLVersion", func() {
		result.HTMLVersion = determineHTMLVersion(rootNode)
	})

	doc := goquery.NewDocumentFromNode(rootNode)
	result.Title = doc.Find("title").Text()

	trackTime("analyzeHeadings", func() { analyzeHeadings(doc, result) })

	// analyze inaccessible links
	analyzeLinks(doc, targetURL, result)

	result.Links.Timeout = linkTimeout.String()
	result.HasLoginForm = checkForLoginForm(doc)

	return result, nil
}

func trackTime(data string, fn func()) {
	start := time.Now()
	fn()
	logger.WithField("duration", time.Since(start)).Info(data)
}

func determineHTMLVersion(rootNode *html.Node) string {
	for c := rootNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.DoctypeNode {
			publicID := ""
			systemID := ""

			for _, attr := range c.Attr {
				switch strings.ToLower(attr.Key) {
				case "public":
					publicID = attr.Val
				case "system":
					systemID = attr.Val
				}
			}

			switch {
			case strings.Contains(publicID, "HTML 4.01"):
				return "HTML 4.01"
			case strings.Contains(publicID, "XHTML 1.0"):
				return "XHTML 1.0"
			case strings.Contains(publicID, "XHTML 1.1"):
				return "XHTML 1.1"
			case strings.TrimSpace(strings.ToLower(c.Data)) == "html":
				return "HTML5"
			default:
				return fmt.Sprintf("Unknown HTML version (publicID: %q, systemID: %q)", publicID, systemID)
			}
		}
	}

	return "HTML5 (no DOCTYPE found, assuming HTML5)"
}

func analyzeHeadings(doc *goquery.Document, result *model.AnalysisResult) {
	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		count := doc.Find(tag).Length()
		result.Headings[tag] = count
	}
}

func analyzeLinks(doc *goquery.Document, baseURL string, result *model.AnalysisResult) {
	start := time.Now()
	defer logger.WithField("duration", time.Since(start)).Info("analyzeLinks")

	base, err := url.Parse(baseURL)
	if err != nil {
		return
	}

	links := doc.Find("a[href]")
	result.Links.TotalLinks = links.Length()

	inAccChan := make(chan string, links.Length())
	var wg sync.WaitGroup

	links.Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		parsedURL, err := url.Parse(href)
		if err != nil {
			return
		}

		// Resolve relative URLs
		resolvedURL := base.ResolveReference(parsedURL)

		if isInternalLink(base, resolvedURL) {
			result.Links.InternalLinks++
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				if !isLinkAccessible(url) {
					inAccChan <- url
				}
			}(resolvedURL.String())
		} else {
			result.Links.ExternalLinks++
		}
	})

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(inAccChan)
	}()

	// Collect results from the channel
	for inacUrl := range inAccChan {
		result.Links.InaccessibleLinks++
		result.Links.InaccessibleLinksList = append(result.Links.InaccessibleLinksList, inacUrl)
	}
}

func isInternalLink(base, target *url.URL) bool {
	return target.Host == "" || target.Host == base.Host
}

func isLinkAccessible(url string) bool {
	client := http.Client{
		Timeout: linkTimeout,
	}

	resp, err := client.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < 400
}

func checkForLoginForm(doc *goquery.Document) bool {
	formPatterns := map[string]string{
		"username": "input[type='text'][name='username'],input[type='email'][name='username'],input[type='text'][name='uname'],input[type='email'][name='uname']",
		"password": "input[type='password'],input[name='password']",
		"submit":   "button:contains('Login'),button:contains('Log In'),button:contains('Sign In')",
	}

	loginFormStatus := false
	doc.Find("form").Each(func(formIndex int, form *goquery.Selection) {
		isUserNameField := form.Find(formPatterns["username"]).Length() > 0
		isPasswordField := form.Find(formPatterns["password"]).Length() > 0
		isSubmitField := form.Find(formPatterns["submit"]).Length() > 0

		if isUserNameField && isPasswordField && isSubmitField {
			loginFormStatus = true
		}
	})

	return loginFormStatus
}
