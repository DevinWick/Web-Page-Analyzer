package handlers

import (
	"github.com/devinwick/web-page-analyzer/model"
	"github.com/devinwick/web-page-analyzer/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

func IndexPathHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func AnalyzeHandler(c *gin.Context) {

	pageUrl := strings.TrimSpace(c.PostForm("url"))

	_, err := url.ParseRequestURI(pageUrl)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": "Invalid URL. Please include a valid URL",
		})
		return
	}

	resp := map[string]any{}

	result, err := service.AnalyzeWebPage(pageUrl)
	if err != nil {
		errorResult := model.AnalysisResult{}
		errorResult.URL = pageUrl
		errorResult.Error = err.Error()
		errorResult.StatusCode = result.StatusCode

		resp["result"] = errorResult
	} else {
		resp["result"] = result
	}

	c.HTML(http.StatusOK, "results.html", resp)
}
