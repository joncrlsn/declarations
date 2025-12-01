package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Declaration represents a single declaration
type Declaration struct {
	ID        int    `json:"id"`
	Label     string `json:"label,omitempty"`
	Text      string `json:"text"`
	Reference string `json:"reference"`
	RawLine   string `json:"raw_line"`
}

// DeclarationsAPI handles all API operations
type DeclarationsAPI struct {
	declarations []Declaration
	nextID       int
	filename     string
	mutex        sync.RWMutex
}

// Bible book order for sorting
var bibleBookOrder = map[string]int{
	"Genesis": 1, "Gen": 1, "Exodus": 2, "Ex": 2, "Exod": 2, "Leviticus": 3, "Lev": 3,
	"Numbers": 4, "Num": 4, "Deuteronomy": 5, "Deut": 5, "Joshua": 6, "Josh": 6,
	"Judges": 7, "Judg": 7, "Ruth": 8, "1 Samuel": 9, "1 Sam": 9, "2 Samuel": 10, "2 Sam": 10,
	"1 Kings": 11, "1 Kgs": 11, "2 Kings": 12, "2 Kgs": 12, "1 Chronicles": 13, "1 Chr": 13,
	"2 Chronicles": 14, "2 Chr": 14, "Ezra": 15, "Nehemiah": 16, "Neh": 16, "Esther": 17, "Est": 17,
	"Job": 18, "Psalms": 19, "Psalm": 19, "Ps": 19, "Proverbs": 20, "Prov": 20,
	"Ecclesiastes": 21, "Eccl": 21, "Song of Songs": 22, "Song": 22, "Isaiah": 23, "Isa": 23,
	"Jeremiah": 24, "Jer": 24, "Lamentations": 25, "Lam": 25, "Ezekiel": 26, "Ezek": 26, "Ez": 26,
	"Daniel": 27, "Dan": 27, "Hosea": 28, "Hos": 28, "Joel": 29, "Amos": 30, "Obadiah": 31, "Obad": 31,
	"Jonah": 32, "Micah": 33, "Mic": 33, "Nahum": 34, "Nah": 34, "Habakkuk": 35, "Hab": 35,
	"Zephaniah": 36, "Zeph": 36, "Haggai": 37, "Hag": 37, "Zechariah": 38, "Zech": 38,
	"Malachi": 39, "Mal": 39, "Matthew": 40, "Matt": 40, "Mark": 41, "Luke": 42, "John": 43,
	"Acts": 44, "Romans": 45, "Rom": 45, "1 Corinthians": 46, "1 Cor": 46, "2 Corinthians": 47, "2 Cor": 47,
	"Galatians": 48, "Gal": 48, "Ephesians": 49, "Eph": 49, "Philippians": 50, "Phil": 50,
	"Colossians": 51, "Col": 51, "1 Thessalonians": 52, "1 Thess": 52, "2 Thessalonians": 53, "2 Thess": 53,
	"1 Timothy": 54, "1 Tim": 54, "2 Timothy": 55, "2 Tim": 55, "Titus": 56, "Philemon": 57,
	"Hebrews": 58, "Heb": 58, "James": 59, "1 Peter": 60, "1 Pet": 60, "2 Peter": 61, "2 Pet": 61,
	"1 John": 62, "2 John": 63, "3 John": 64, "Jude": 65, "Revelation": 66, "Rev": 66,
}

// NewDeclarationsAPI creates a new API instance
func NewDeclarationsAPI(filename string) *DeclarationsAPI {
	api := &DeclarationsAPI{
		declarations: make([]Declaration, 0),
		nextID:       1,
		filename:     filename,
	}

	rand.Seed(time.Now().UnixNano())

	if err := api.loadDeclarations(); err != nil {
		fmt.Printf("Error loading declarations: %v\n", err)
	}

	return api
}

// parseDeclaration parses a line from the file
func (api *DeclarationsAPI) parseDeclaration(line string, id int) Declaration {
	decl := Declaration{
		ID:      id,
		RawLine: line,
	}

	// Check for labels (multiple labels separated by colons)
	// Format: :Label1:Label2:Label3: followed by text
	labelRegex := regexp.MustCompile(`^:(([^:\s]+:)+)\s*(.*)$`)
	if matches := labelRegex.FindStringSubmatch(line); matches != nil {
		// Remove trailing colon from captured labels
		decl.Label = strings.TrimSuffix(matches[1], ":")
		line = matches[3]
	}

	// Split by " - " to separate text and reference
	parts := strings.SplitN(line, " - ", 2)
	if len(parts) == 2 {
		decl.Text = strings.TrimSpace(parts[0])
		decl.Reference = strings.TrimSpace(parts[1])
	} else {
		decl.Text = strings.TrimSpace(line)
	}

	return decl
}

// loadDeclarations loads declarations from file
func (api *DeclarationsAPI) loadDeclarations() error {
	api.mutex.Lock()
	defer api.mutex.Unlock()

	file, err := os.Open(api.filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	api.declarations = make([]Declaration, 0)
	api.nextID = 1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		decl := api.parseDeclaration(line, api.nextID)
		api.declarations = append(api.declarations, decl)
		api.nextID++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

// parseReference extracts book, chapter, and verse for sorting
func parseReference(reference string) (book string, chapter int, verse int) {
	reference = strings.TrimSpace(reference)
	if reference == "" {
		return "", 0, 0
	}

	// Take first reference segment before separators
	primary := reference
	delims := []string{";", " and ", "&", ", "}
	minIdx := len(primary)
	for _, d := range delims {
		if idx := strings.Index(primary, d); idx >= 0 && idx < minIdx {
			minIdx = idx
		}
	}
	if minIdx != len(primary) {
		primary = primary[:minIdx]
	}
	primary = strings.TrimSpace(primary)

	// Check if this is a Bible reference (contains chapter:verse format) or a person name
	parts := strings.Fields(primary)
	if len(parts) < 2 {
		return primary, 0, 0
	}

	// Check if the last part contains a colon (indicating chapter:verse)
	lastPart := parts[len(parts)-1]
	if !strings.Contains(lastPart, ":") {
		// No colon found - this is likely a person name, return full reference
		return primary, 0, 0
	}

	// Parse "Book Chapter:Verse"
	chapterVerse := lastPart
	book = strings.Join(parts[:len(parts)-1], " ")

	cvParts := strings.Split(chapterVerse, ":")
	if len(cvParts) >= 2 {
		if c, err := strconv.Atoi(cvParts[0]); err == nil {
			chapter = c
		}
		// Handle verse ranges like "14-15"
		versePart := strings.Split(cvParts[1], "-")[0]
		versePart = strings.Split(versePart, ",")[0]
		if v, err := strconv.Atoi(versePart); err == nil {
			verse = v
		}
	}

	return book, chapter, verse
}

// sortDeclarations sorts by Bible book order
func (api *DeclarationsAPI) sortDeclarations() {
	sort.Slice(api.declarations, func(i, j int) bool {
		bookI, chapterI, verseI := parseReference(api.declarations[i].Reference)
		bookJ, chapterJ, verseJ := parseReference(api.declarations[j].Reference)

		orderI, existsI := bibleBookOrder[bookI]
		orderJ, existsJ := bibleBookOrder[bookJ]

		// Both are Bible books
		if existsI && existsJ {
			if orderI != orderJ {
				return orderI < orderJ
			}
			if chapterI != chapterJ {
				return chapterI < chapterJ
			}
			return verseI < verseJ
		}

		// Bible books come before non-Bible references
		if existsI && !existsJ {
			return true
		}
		if !existsI && existsJ {
			return false
		}

		// Both are non-Bible, sort alphabetically
		return bookI < bookJ
	})
}

// saveDeclarations saves declarations to file
func (api *DeclarationsAPI) saveDeclarations() error {
	api.sortDeclarations()

	file, err := os.Create(api.filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	for _, decl := range api.declarations {
		var line string
		if decl.Label != "" {
			line = fmt.Sprintf(":%s: %s - %s", decl.Label, decl.Text, decl.Reference)
		} else {
			line = fmt.Sprintf("%s - %s", decl.Text, decl.Reference)
		}

		if _, err := fmt.Fprintln(file, line); err != nil {
			return fmt.Errorf("failed to write line: %w", err)
		}
	}

	return nil
}

// SortAndSave sorts and saves declarations (for SORT_ONLY mode)
func (api *DeclarationsAPI) SortAndSave() error {
	api.mutex.Lock()
	defer api.mutex.Unlock()
	return api.saveDeclarations()
}

// findDeclarationByID finds a declaration by ID
func (api *DeclarationsAPI) findDeclarationByID(id int) (*Declaration, int) {
	for i, decl := range api.declarations {
		if decl.ID == id {
			return &api.declarations[i], i
		}
	}
	return nil, -1
}

// API Handlers

// GetDeclarations handles GET /api/v1/declarations
func (api *DeclarationsAPI) GetDeclarations(w http.ResponseWriter, r *http.Request) {
	api.mutex.RLock()
	defer api.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"declarations": api.declarations,
	})
}

// GetDeclaration handles GET /api/v1/declarations/{id}
func (api *DeclarationsAPI) GetDeclaration(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	api.mutex.RLock()
	defer api.mutex.RUnlock()

	decl, _ := api.findDeclarationByID(id)
	if decl == nil {
		http.Error(w, "Declaration not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(decl)
}

// GetRandomDeclaration handles GET /api/v1/declarations/random
func (api *DeclarationsAPI) GetRandomDeclaration(w http.ResponseWriter, r *http.Request) {
	api.mutex.RLock()
	defer api.mutex.RUnlock()

	if len(api.declarations) == 0 {
		http.Error(w, "No declarations available", http.StatusNotFound)
		return
	}

	randomIndex := rand.Intn(len(api.declarations))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.declarations[randomIndex])
}

// GetDeclarationsByLabel handles GET /api/v1/declarations/label/{label}
func (api *DeclarationsAPI) GetDeclarationsByLabel(w http.ResponseWriter, r *http.Request, label string) {
	api.mutex.RLock()
	defer api.mutex.RUnlock()

	var filtered []Declaration
	for _, decl := range api.declarations {
		if decl.Label != "" {
			// Check if the label matches (handle multiple labels separated by colons)
			labels := strings.Split(decl.Label, ":")
			for _, l := range labels {
				if strings.TrimSpace(l) == label {
					filtered = append(filtered, decl)
					break
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"declarations": filtered,
		"label":        label,
	})
}

// LabelCount represents a label with its declaration count
type LabelCount struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

// GetLabels handles GET /api/v1/labels
func (api *DeclarationsAPI) GetLabels(w http.ResponseWriter, r *http.Request) {
	api.mutex.RLock()
	defer api.mutex.RUnlock()

	// Count declarations per label
	labelCounts := make(map[string]int)
	for _, decl := range api.declarations {
		if decl.Label != "" {
			labels := strings.Split(decl.Label, ":")
			for _, l := range labels {
				label := strings.TrimSpace(l)
				if label != "" {
					labelCounts[label]++
				}
			}
		}
	}

	// Convert to sorted slice
	var labels []LabelCount
	for label, count := range labelCounts {
		labels = append(labels, LabelCount{Label: label, Count: count})
	}

	// Sort alphabetically
	sort.Slice(labels, func(i, j int) bool {
		return labels[i].Label < labels[j].Label
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"labels": labels,
	})
}

// Note: Create, Update, and Delete operations removed - API is now read-only

// GetHealth handles GET /api/v1/health
func (api *DeclarationsAPI) GetHealth(w http.ResponseWriter, r *http.Request) {
	api.mutex.RLock()
	defer api.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":             "healthy",
		"declarations_count": len(api.declarations),
	})
}

// GetEnv handles GET /api/v1/env
func GetEnv(w http.ResponseWriter, r *http.Request) {
	envVars := os.Environ()
	envMap := make(map[string]string)

	for _, env := range envVars {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"environment": envMap,
	})
}

// getESVToken retrieves the ESV API token from environment or file
func getESVToken() (string, error) {
	// Check environment variable first (Populated for Google Cloud Run)
	if token := os.Getenv("ESV_API_TOKEN"); token != "" {
		return strings.TrimSpace(token), nil
	}

	// Fall back to local file
	tokenData, err := os.ReadFile(".esv-api-token")
	if err != nil {
		return "", fmt.Errorf("failed to read API token: %w", err)
	}
	return strings.TrimSpace(string(tokenData)), nil
}

// GetBibleText handles GET /api/v1/bible-esv/{reference}
// Proxies requests to ESV API to protect the token
func GetBibleText(w http.ResponseWriter, r *http.Request, reference string) {
	if reference == "" {
		http.Error(w, "Missing reference", http.StatusBadRequest)
		return
	}

	// Get API token
	token, err := getESVToken()
	if err != nil {
		http.Error(w, "Failed to get API token", http.StatusInternalServerError)
		return
	}

	// Make request to ESV API with properly encoded query parameter
	esvURL := "https://api.esv.org/v3/passage/text/?q=" + url.QueryEscape(reference)
	req, err := http.NewRequest("GET", esvURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Token "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to call ESV API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
