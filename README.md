# Web-Page-Analyzer
This is a web application that allows you to analyze web page. You can submit a web page url and web application will return below information

- HTML version
- Page Title
- Heading count by heading level
- Internal Link Count
- External Link Count
- Inaccessible Link Count
- Inaccessible Link Details
- Whether web page contains a login form

---
## Content
1. Technologies
2. Prerequisites
3. Installation
4. Usage
5. API Endpoints
6. Challenges & Solutions
7. Future Improvements
8. Assumptions

## Technologies

### Backend
- Go
- Gin Web Framework https://github.com/gin-gonic/gin
- Go Query(html parsing) https://github.com/PuerkitoBio/goquery
- Logrus (logging) https://github.com/sirupsen/logrus
- Testify (Testing) 

### Frontend
- HTML5, CSS3
- Gin HTML Templates

### DevOps
- Go Modules (Dependency Management)
- GO Testing Framework

## Prerequisites
- Go (version 1.21 higher)
- Git

## Installation
1. Unzip the file and open the directory
```
cd Web-Page-Analyzer
```
2. Install Dependencies
```
go mod download
```

3. Build and Run


 - Build the project and run the executable
```
go build -o web-analyzer
./web-analyzer
```

### Usage
1. Access the application at ```http://localhost:8080/```
2. Enter valid URL in the field
3. Click "Analyze" to submit
4. View the analysis result
5. You can view the logs in the console

Please note that if links in the web page are not responded withing **timeout period (10s)** that link will be counted as **inaccessible link**

### API Endpoint
Below endpoint is provided to analyze the webpage
``` 
POST /analyze
Content-Type: application/x-www-form-urlencoded
```
Response
- response contain html page with below information
``` 
Site Details
    -URL
    -Title
    -HTML Version
    -Is Containg a Login Form
Headings Count
Links Count
    -Internal Links
    -External Links
    -Inaccessible Links
    -Inaccessible Links Details
```
### Challanges and Solutions
1. HTML Parsing & Anlaysis
- Challenge: Accurately determining the doctype declarations of old html web pages
- Solution: Implemented multi-level checking from DOCTYPE node and public attributes
2. InAccessible and Accessible Link Checking
 - Challenge: Sequential checking was slow for pages with too many links
 - Solution: Implemented concurrent calls for link checking
3. Performance Optimization
- Challenge: Some links timeouts after taking a long period
- Solution: timeouts the long-lasting requests in the predefined timeout period and count as inaccessible link
4. Bot mitigation Tools
- Challenges: Some web pages are not reachable due to bot mitigation techniques used
- Solution: need to develop a solution to emulate human like behaviour as in browser 
- Handling cookies and required request headers looks by specific bot mitigation solution

### Future Improvements
1. Implement caching of analyzed results for better response time and minimize resource utlization
2. Add rate limiting to prevent attacks or abused use
3. Record analysis in database for history data analysis
4. Add performance metrics of page in analysis (Page Load Time, SEO)
5. Containerize and make timeout periods configurable through env variables
6. Add Progress indicator for long analysis

### Assumptions
1. HTML Version Detection
 - Assumption: html version analyzed only using DOCTYPE declaration
 - Only checks for common DOCTYPE patterns (HTML5, HTML4, XHTML)
2. Login Form Detection
- Assumes that all login forms contains ```form``` tag
- Assumes that username, password and submit button must be inside a ```form``` tag to be a login form
3. Link Analysis
- Assumes HEAD requests are sufficient
- Timeout of 10s per link check
4. Error Handling
- Assumes pages return valid HTML page
- Assumes valid url must contain http or https protocol
5. Security
- Assumes URLs are safe to fetch and follows http redirects by default
- No rate limiting implemented.
