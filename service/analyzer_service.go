package service

import (
	"errors"
	"github.com/devinwick/web-page-analyzer/model"
)

func AnalyzeWebPage(targetURL string) (*model.AnalysisResult, error) {
	result := &model.AnalysisResult{
		URL:      targetURL,
		Headings: make(map[string]int),
	}

	result.StatusCode = 200

	return result, errors.New("test error")
}
