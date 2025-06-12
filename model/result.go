package model

type AnalysisResult struct {
	StatusCode   int
	URL          string
	Title        string
	HTMLVersion  string
	Headings     map[string]int
	Links        LinkAnalysis
	HasLoginForm bool
	Error        string
}

type LinkAnalysis struct {
	InternalLinks         int
	ExternalLinks         int
	InaccessibleLinks     int
	TotalLinks            int
	InaccessibleLinksList []string
	Timeout               string
}
