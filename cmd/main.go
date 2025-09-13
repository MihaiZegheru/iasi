package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// main is the entry point for the CLI tool. It fetches, filters, groups, sorts, and writes the user's 100-point problems to CSV.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: iasi {username}")
		os.Exit(1)
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

	filtered := filter100PointEntries(records)
	grouped := groupByProblemEarliest(filtered)
	final := sortByDate(grouped)

	outDir := "data"
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		if err := os.Mkdir(outDir, 0755); err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}
	outPath := outDir + string(os.PathSeparator) + username + "_timeline.csv"
	if err := writeCSV(outPath, final); err != nil {
		log.Fatalf("Failed to write CSV: %v", err)
	}
	fmt.Printf("Saved %d entries to %s\n", len(final), outPath)
}

// fetchAllEntries paginates and fetches all monitor entries for a given username.
func fetchAllEntries(username string) ([][]string, error) {
	var records [][]string
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
			s.Find("td").Each(func(j int, td *goquery.Selection) {
				row = append(row, strings.TrimSpace(td.Text()))
			})
			if len(row) > 0 {
				records = append(records, row)
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
func writeCSV(filename string, records [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"name", "url", "time"})

	for _, record := range records {
		if len(record) >= 6 {
			id := record[0]
			if strings.HasPrefix(id, "#") {
				id = id[1:]
			}
			url := "https://www.infoarena.ro/job_detail/" + id
			pruned := []string{record[2], url, record[5]}
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
