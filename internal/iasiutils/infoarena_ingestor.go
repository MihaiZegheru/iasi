package iasiutils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// InfoarenaIngestor handles fetching and parsing Infoarena problems and solutions
// All logic for scraping Infoarena should go here.
type InfoarenaIngestor struct{}

func (ii *InfoarenaIngestor) FetchProblemAndSolution(id string) (string, string, error) {
	// 1. Fetch the job_detail page for the solution (for problem link)
	jobURL := "https://www.infoarena.ro/job_detail/" + id
	log.Printf("[DEBUG] Fetching job_detail page: %s", jobURL)
	resp, err := http.Get(jobURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)
	log.Printf("[DEBUG] job_detail HTML (first 500 chars): %s", TruncateString(bodyStr, 500))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
	if err != nil {
		return "", "", err
	}
	// 2. Find the link to the problem page
	problemURL := ""
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.HasPrefix(href, "/problema/") {
			problemURL = "https://www.infoarena.ro" + href
		}
	})
	log.Printf("[DEBUG] Extracted problemURL: %s", problemURL)
	if problemURL == "" {
		return "", "", fmt.Errorf("problem URL not found on job_detail page")
	}
	// 3. Fetch the problem page for the statement
	log.Printf("[DEBUG] Fetching problem page: %s", problemURL)
	resp2, err := http.Get(problemURL)
	if err != nil {
		return "", "", err
	}
	defer resp2.Body.Close()
	body2Bytes, _ := ioutil.ReadAll(resp2.Body)
	body2Str := string(body2Bytes)
	log.Printf("[DEBUG] problem page HTML (first 500 chars): %s", TruncateString(body2Str, 500))
	doc2, err := goquery.NewDocumentFromReader(strings.NewReader(body2Str))
	if err != nil {
		return "", "", err
	}
	// 4. Extract the problem statement (main content)
	statement := strings.TrimSpace(doc2.Find(".wiki_text_block").Text())
	if statement == "" {
		// fallback: try .content .problem-text
		statement = strings.TrimSpace(doc2.Find(".content .problem-text").Text())
	}
	if statement == "" {
		// fallback: try just .content
		statement = strings.TrimSpace(doc2.Find(".content").Text())
	}
	if statement == "" {
		// fallback: try body text
		statement = strings.TrimSpace(doc2.Find("body").Text())
	}
	log.Printf("[DEBUG] Extracted statement (first 200 chars): %s", TruncateString(statement, 200))

	// 5. Fetch the solution from job_detail/{id}?action=view-source
	solutionURL := jobURL + "?action=view-source"
	log.Printf("[DEBUG] Fetching solution page: %s", solutionURL)
	resp3, err := http.Get(solutionURL)
	if err != nil {
		return statement, "", err
	}
	defer resp3.Body.Close()
	solutionBytes, _ := ioutil.ReadAll(resp3.Body)
	solutionStr := string(solutionBytes)
	log.Printf("[DEBUG] solution page HTML (first 500 chars): %s", TruncateString(solutionStr, 500))
	doc3, err := goquery.NewDocumentFromReader(strings.NewReader(solutionStr))
	if err != nil {
		return statement, "", err
	}

	// Check if the force_view_source form/button is present
	if doc3.Find("#force_view_source").Length() > 0 {
		log.Printf("[INFO] 'Vezi sursa' button detected. Submitting form to reveal source code.")
		client := &http.Client{Timeout: 30 * time.Second}
		formData := "force_view_source=Vezi+sursa"
		req, err := http.NewRequest("POST", solutionURL, strings.NewReader(formData))
		if err != nil {
			return statement, "", err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp4, err := client.Do(req)
		if err != nil {
			return statement, "", err
		}
		defer resp4.Body.Close()
		solutionBytes, _ = ioutil.ReadAll(resp4.Body)
		solutionStr = string(solutionBytes)
	log.Printf("[DEBUG] solution page after form submit (first 500 chars): %s", TruncateString(solutionStr, 500))
		doc3, err = goquery.NewDocumentFromReader(strings.NewReader(solutionStr))
		if err != nil {
			return statement, "", err
		}
	}

	// Infoarena solution code may be nested in <code class="hljs cpp"> with inner spans, so concatenate all text nodes
	// Try all <code>, <pre>, <textarea> tags in order, recursively extracting all text
	extractAllText := func(sel *goquery.Selection) string {
		var sb strings.Builder
		var extract func(*goquery.Selection)
		extract = func(s *goquery.Selection) {
			s.Contents().Each(func(i int, n *goquery.Selection) {
				if goquery.NodeName(n) == "#text" {
					sb.WriteString(n.Text())
				} else {
					extract(n)
				}
			})
		}
		extract(sel)
		return sb.String()
	}
	var solutionBuilder strings.Builder
	doc3.Find("code").Each(func(i int, sel *goquery.Selection) {
		solutionBuilder.WriteString(extractAllText(sel))
		solutionBuilder.WriteString("\n")
	})
	doc3.Find("pre").Each(func(i int, sel *goquery.Selection) {
		solutionBuilder.WriteString(extractAllText(sel))
		solutionBuilder.WriteString("\n")
	})
	doc3.Find("textarea").Each(func(i int, sel *goquery.Selection) {
		solutionBuilder.WriteString(extractAllText(sel))
		solutionBuilder.WriteString("\n")
	})
	solution := strings.TrimSpace(solutionBuilder.String())
	log.Printf("[DEBUG] Extracted solution (first 200 chars): %s", TruncateString(solution, 200))
	return statement, solution, nil
}
