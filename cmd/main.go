package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// buildLLMPrompt creates a prompt for the LLM using the problem statement and solution
func buildLLMPrompt(statement, solution string) string {
       // Defensive: if statement or solution is empty, say so in the prompt
       if strings.TrimSpace(statement) == "" {
	       statement = "(Problem statement could not be fetched)"
       }
       if strings.TrimSpace(solution) == "" {
	       solution = "(Solution code could not be fetched)"
       }
       return fmt.Sprintf(`You are an expert competitive programming assistant. Given the following problem statement and its solution, generate:
       - 3 helpful hints for a student (in Romanian, do not give away the full solution)
       - a detailed editorial (in Romanian, explaining the solution and key ideas)

       Problem statement:
       %s

       Solution:
       %s

       Return a JSON object with two fields: "hints" (an array of 3 strings) and "editorial" (a string).`, statement, solution)
}

// fetchInfoarenaProblem fetches the problem statement and solution text from Infoarena for a given job_detail id
func fetchInfoarenaProblem(id string) (string, string, error) {
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
	log.Printf("[DEBUG] job_detail HTML (first 500 chars): %s", truncateString(bodyStr, 500))
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
	log.Printf("[DEBUG] problem page HTML (first 500 chars): %s", truncateString(body2Str, 500))
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
       log.Printf("[DEBUG] Extracted statement (first 200 chars): %s", truncateString(statement, 200))

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
	log.Printf("[DEBUG] solution page HTML (first 500 chars): %s", truncateString(solutionStr, 500))
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
		log.Printf("[DEBUG] solution page after form submit (first 500 chars): %s", truncateString(solutionStr, 500))
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
	log.Printf("[DEBUG] Extracted solution (first 200 chars): %s", truncateString(solution, 200))
	return statement, solution, nil
	}

// truncateString returns the first n characters of s, appending ... if truncated
func truncateString(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
// callGeminiLLM calls the Gemini LLM API with the prompt and returns the response JSON
func callGeminiLLM(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		 return "", fmt.Errorf("GEMINI_API_KEY not set")
	}
	url := "https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key=" + apiKey
	reqBody := fmt.Sprintf(`{"contents":[{"parts":[{"text":%q}]}]}`, prompt)
	req, err := http.NewRequest("POST", url, strings.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// Log the raw Gemini response for debugging
	log.Printf("[DEBUG] Raw Gemini API response: %s", truncateString(string(body), 1000))
	// Parse Gemini response
	var parsed struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}
	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("No LLM response candidates. Raw response: %s", truncateString(string(body), 1000))
	}
	return parsed.Candidates[0].Content.Parts[0].Text, nil
}

// main is the entry point for the CLI tool. It fetches, filters, groups, sorts, and writes the user's 100-point problems to CSV.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: iasi <username> OR iasi run <username>")
		os.Exit(1)
	}
	if os.Args[1] == "run" && len(os.Args) >= 3 {
		username := os.Args[2]
		serveTracker(username)
		return
	}
	username := os.Args[1]

	records, err := fetchAllEntries(username)
	if err != nil {
		log.Fatalf("Error fetching entries: %v", err)
	}
	if len(records) == 0 {
		fmt.Println("No entries found for user.")
		return
	}

	// Convert []monitorRow to [][]string for filter/group/sort
	var raw [][]string
	for _, r := range records {
		raw = append(raw, r.fields)
	}
	filtered := filter100PointEntries(raw)
	grouped := groupByProblemEarliest(filtered)
	final := sortByDate(grouped)

	// Map back to monitorRow to get problemUrl
	var finalRows []monitorRow
	for _, row := range final {
		for _, r := range records {
			if len(row) == len(r.fields) && row[0] == r.fields[0] && row[2] == r.fields[2] && row[5] == r.fields[5] {
				finalRows = append(finalRows, r)
				break
			}
		}
	}

	outDir := "data"
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		if err := os.Mkdir(outDir, 0755); err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}
	outPath := outDir + string(os.PathSeparator) + username + "_timeline.csv"
	if err := writeCSV(outPath, finalRows); err != nil {
		log.Fatalf("Failed to write CSV: %v", err)
	}
	fmt.Printf("Saved %d entries to %s\n", len(finalRows), outPath)
}
// serveTracker starts a web server to show the tracker UI and serve the problem list as JSON.
func serveTracker(username string) {
	// Log to console only (no debug.log file)
	log.SetOutput(os.Stdout)
	log.Println("[INFO] serveTracker started for user:", username)

	// From here, only Go logs go to debug.log. React dev server output goes to console.

	// --- LLM Editorial/Hints API ---
	// POST /problems/{id}/generate
	http.HandleFunc("/problems/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/problems/")
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			log.Printf("[ERROR] Invalid /problems/ path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		id := parts[0]
		action := parts[1]
		editorialPath := "data/editorials/" + id + ".json"
		if action == "generate" && r.Method == "POST" {
			log.Printf("[INFO] /problems/%s/generate POST called", id)
			// Check cache first
			if _, err := os.Stat(editorialPath); err == nil {
				log.Printf("[INFO] Editorial cache hit for %s", editorialPath)
				data, _ := ioutil.ReadFile(editorialPath)
				w.Header().Set("Content-Type", "application/json")
				w.Write(data)
				return
			}
		       log.Printf("[INFO] Fetching problem and solution for id %s", id)
		       statement, solution, err := fetchInfoarenaProblem(id)
		       if err != nil {
			       log.Printf("[ERROR] Failed to fetch problem/solution: %v", err)
			       http.Error(w, "Failed to fetch problem/solution: "+err.Error(), 500)
			       return
		       }
		       if strings.TrimSpace(statement) == "" || strings.TrimSpace(solution) == "" {
			       log.Printf("[ERROR] Statement or solution missing. Statement: '%s' Solution: '%s'", truncateString(statement, 100), truncateString(solution, 100))
			       http.Error(w, "Problem statement or solution could not be fetched. Please check the Infoarena page structure.", 500)
			       return
		       }
		       log.Printf("[INFO] Problem and solution fetched. Building prompt.")
		       prompt := buildLLMPrompt(statement, solution)
		       log.Printf("[DEBUG] Prompt: %s", prompt)
		       llmResp, err := callGeminiLLM(prompt)
		       if err != nil {
			       log.Printf("[ERROR] LLM error: %v", err)
			       http.Error(w, "LLM error: "+err.Error(), 500)
			       return
		       }
		       log.Printf("[INFO] LLM response received. Raw response: %s", llmResp)
		       log.Printf("[INFO] Attempting to parse JSON.")
	       var result map[string]interface{}
	       llmJson := llmResp
	       // Try to extract JSON from code block or text if direct parse fails
	       if err := json.Unmarshal([]byte(llmJson), &result); err != nil {
		       log.Printf("[WARN] Direct JSON parse failed: %v", err)
		       // Try to extract JSON from markdown/code block or text
		       jsonStart := strings.Index(llmResp, "{")
		       jsonEnd := strings.LastIndex(llmResp, "}")
		       if jsonStart != -1 && jsonEnd > jsonStart {
			       llmJson = llmResp[jsonStart : jsonEnd+1]
			       if err2 := json.Unmarshal([]byte(llmJson), &result); err2 == nil {
				       log.Printf("[INFO] JSON extracted from LLM output.")
			       } else {
				       log.Printf("[WARN] JSON extraction also failed: %v", err2)
				       result = map[string]interface{}{
					       "hints": []string{"LLM output could not be parsed as JSON."},
					       "editorial": llmResp,
				       }
			       }
		       } else {
			       result = map[string]interface{}{
				       "hints": []string{"LLM output could not be parsed as JSON."},
				       "editorial": llmResp,
			       }
		       }
	       }
	       jsonBytes, _ := json.MarshalIndent(result, "", "  ")
	       ioutil.WriteFile(editorialPath, jsonBytes, 0644)
	       w.Header().Set("Content-Type", "application/json")
	       w.Write(jsonBytes)
	       log.Printf("[INFO] Editorial for %s generated and returned.", id)
	       return
		}
		if action == "editorial" && r.Method == "GET" {
			if _, err := os.Stat(editorialPath); err == nil {
				log.Printf("[INFO] Editorial cache GET for %s", editorialPath)
				data, _ := ioutil.ReadFile(editorialPath)
				w.Header().Set("Content-Type", "application/json")
				w.Write(data)
				return
			}
			log.Printf("[WARN] Editorial not generated for %s", editorialPath)
			http.Error(w, "Not generated", http.StatusNotFound)
			return
		}
		log.Printf("[ERROR] Unknown /problems/ action: %s", action)
		http.NotFound(w, r)
	})
	// Start React dev server (output to console)
	var reactCmd *exec.Cmd
	if os.PathSeparator == '\\' { // Windows
		reactCmd = exec.Command("cmd", "/C", "cd web/tracker-app && npm run dev")
	} else {
		reactCmd = exec.Command("sh", "-c", "cd web/tracker-app && npm run dev")
	}
	reactCmd.Stdout = os.Stdout
	reactCmd.Stderr = os.Stderr
	if err := reactCmd.Start(); err != nil {
		log.Fatalf("Failed to start React dev server: %v", err)
	}

	// Wait for React dev server to be ready
	ready := false
	for i := 0; i < 30; i++ {
		time.Sleep(1 * time.Second)
		resp, err := http.Get("http://localhost:5173")
		if err == nil && resp.StatusCode == 200 {
			ready = true
			resp.Body.Close()
			break
		}
	}
	if !ready {
		log.Println("Warning: React dev server did not become ready in time.")
	}

	// Open browser to React app
	openBrowser("http://localhost:5173/")

	// On exit, kill React dev server
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		_ = reactCmd.Process.Kill()
		os.Exit(0)
	}()

	records, err := fetchAllEntries(username)
	if err != nil {
		log.Fatalf("Error fetching entries: %v", err)
	}
	// Convert []monitorRow to [][]string for filter/group/sort
	var raw [][]string
	for _, r := range records {
		raw = append(raw, r.fields)
	}
	filtered := filter100PointEntries(raw)
	grouped := groupByProblemEarliest(filtered)
	final := sortByDate(grouped)

	// Map back to monitorRow to get problemUrl
	var finalRows []monitorRow
	for _, row := range final {
		for _, r := range records {
			if len(row) == len(r.fields) && row[0] == r.fields[0] && row[2] == r.fields[2] && row[5] == r.fields[5] {
				finalRows = append(finalRows, r)
				break
			}
		}
	}

	type Problem struct {
		Name        string `json:"name"`
		Url         string `json:"url"`
		UrlSolution string `json:"url_solution"`
		Time        string `json:"time"`
		Id          string `json:"id"`
	}
	var problems []Problem
	for _, r := range finalRows {
		fields := r.fields
		if len(fields) >= 6 {
			// Try to extract job_detail id from problemUrl if possible
			id := ""
			if strings.Contains(r.problemUrl, "/job_detail/") {
				parts := strings.Split(r.problemUrl, "/job_detail/")
				if len(parts) > 1 {
					id = parts[1]
				}
			} else {
				id = fields[0]
				if strings.HasPrefix(id, "#") {
					id = id[1:]
				}
			}
			urlSolution := "https://www.infoarena.ro/job_detail/" + id
			problems = append(problems, Problem{
				Name:        fields[2],
				Url:         r.problemUrl,
				UrlSolution: urlSolution,
				Time:        fields[5],
				Id:          id,
			})
		}
	}

	http.HandleFunc("/problems", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"username":%q,"problems":`, username)
		fmt.Fprint(w, "[")
		for i, p := range problems {
			if i > 0 {
				fmt.Fprint(w, ",")
			}
			fmt.Fprintf(w, `{"name":%q,"url":%q,"time":%q,"id":%q}`, p.Name, p.Url, p.Time, p.Id)
		}
		fmt.Fprint(w, "]}")
	})

		// On exit, kill React dev server (disabled for debugging)
		// go func() {
		// 	ch := make(chan os.Signal, 1)
		// 	signal.Notify(ch, os.Interrupt)
		// 	<-ch
		// 	_ = reactCmd.Process.Kill()
		// 	os.Exit(0)
		// }()

	log.Println("Go API server running at http://localhost:8080 (API only, UI at http://localhost:5173)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// openBrowser tries to open the URL in the default browser (Windows only for now).
func openBrowser(url string) {
	execCmd := "start " + url
	_ = execCommand(execCmd)
}

// execCommand runs a shell command (Windows only).
func execCommand(cmd string) error {
	return exec.Command("cmd", "/C", cmd).Start()
}

// fetchAllEntries paginates and fetches all monitor entries for a given username.
type monitorRow struct {
	fields     []string
	problemUrl string
}

func fetchAllEntries(username string) ([]monitorRow, error) {
	var records []monitorRow
	pageSize := 250
	for offset := 0; ; offset += pageSize {
		url := fmt.Sprintf("https://www.infoarena.ro/monitor?user=%s&display_entries=%d&first_entry=%d", username, pageSize, offset)
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch URL: %w", err)
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to parse HTML: %w", err)
		}
		entriesOnPage := 0
		doc.Find("table.monitor tbody tr").Each(func(i int, s *goquery.Selection) {
			var row []string
			problemUrl := ""
			s.Find("td").Each(func(j int, td *goquery.Selection) {
				row = append(row, strings.TrimSpace(td.Text()))
				if j == 2 { // 3rd column: problem name and link
					if a := td.Find("a"); a.Length() > 0 {
						href, exists := a.Attr("href")
						if exists && strings.HasPrefix(href, "/problema/") {
							problemUrl = "https://www.infoarena.ro" + href
						}
					}
				}
			})
			if len(row) > 0 {
				records = append(records, monitorRow{fields: row, problemUrl: problemUrl})
				entriesOnPage++
			}
		})
		if entriesOnPage < pageSize {
			break // Last page reached
		}
	}
	return records, nil
}

// filter100PointEntries returns only entries with exactly 'Evaluare completa: 100 puncte' in the last column.
func filter100PointEntries(records [][]string) [][]string {
	var filtered [][]string
	for _, record := range records {
		if len(record) > 0 && strings.TrimSpace(record[len(record)-1]) == "Evaluare completa: 100 puncte" {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

// groupByProblemEarliest groups by problem name (3rd column) and keeps the earliest submission (by 6th column, date).
func groupByProblemEarliest(records [][]string) map[string][]string {
	grouped := make(map[string][]string)
	for _, record := range records {
		if len(record) < 6 {
			continue
		}
		name := record[2]
		date := record[5]
		if prev, ok := grouped[name]; ok {
			if compareInfoarenaDate(date, prev[5]) < 0 {
				grouped[name] = record
			}
		} else {
			grouped[name] = record
		}
	}
	return grouped
}

// sortByDate sorts the grouped map values by date (6th column), ascending.
func sortByDate(grouped map[string][]string) [][]string {
	var final [][]string
	for _, v := range grouped {
		final = append(final, v)
	}
	if len(final) > 1 {
		type rowWithTime struct {
			row []string
			t   int64
		}
		var rows []rowWithTime
		for _, r := range final {
			tm, err := parseInfoarenaDate(r[5])
			var t int64
			if err == nil {
				t = tm.Unix()
			}
			rows = append(rows, rowWithTime{row: r, t: t})
		}
		// Sort by t
		for i := 0; i < len(rows)-1; i++ {
			for j := 0; j < len(rows)-i-1; j++ {
				if rows[j].t > rows[j+1].t {
					rows[j], rows[j+1] = rows[j+1], rows[j]
				}
			}
		}
		final = nil
		for _, r := range rows {
			final = append(final, r.row)
		}
	}
	return final
}

// writeCSV writes the records to a CSV file.
func writeCSV(filename string, records []monitorRow) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"name", "url", "url_solution", "time"})

	for _, r := range records {
		fields := r.fields
		if len(fields) >= 6 {
			id := fields[0]
			if strings.HasPrefix(id, "#") {
				id = id[1:]
			}
			urlSolution := "https://www.infoarena.ro/job_detail/" + id
			pruned := []string{fields[2], r.problemUrl, urlSolution, fields[5]}
			writer.Write(pruned)
		}
	}
	return nil
}

// parseInfoarenaDate parses dates like "1 apr 25 13:06:35" to time.Time.
func parseInfoarenaDate(s string) (time.Time, error) {
	months := map[string]string{"ian": "Jan", "feb": "Feb", "mar": "Mar", "apr": "Apr", "mai": "May", "iun": "Jun", "iul": "Jul", "aug": "Aug", "sep": "Sep", "oct": "Oct", "nov": "Nov", "dec": "Dec"}
	parts := strings.Split(s, " ")
	if len(parts) != 4 {
		return time.Time{}, fmt.Errorf("invalid date format")
	}
	day := parts[0]
	mon, ok := months[parts[1]]
	if !ok {
		return time.Time{}, fmt.Errorf("invalid month")
	}
	year := parts[2]
	timepart := parts[3]
	dateStr := fmt.Sprintf("%s %s %s %s", day, mon, year, timepart)
	return time.Parse("2 Jan 06 15:04:05", dateStr)
}

// compareInfoarenaDate compares two infoarena date strings. Returns -1 if a < b, 1 if a > b, 0 if equal or error.
func compareInfoarenaDate(a, b string) int {
	ta, ea := parseInfoarenaDate(a)
	tb, eb := parseInfoarenaDate(b)
	if ea != nil || eb != nil {
		return 0
	}
	if ta.Before(tb) {
		return -1
	}
	if ta.After(tb) {
		return 1
	}
	return 0
}
