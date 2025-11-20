package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Declaration represents a single declaration with its metadata
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
	mutex        sync.RWMutex
	filename     string
}

// Bible book order for sorting declarations
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

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	if err := api.loadDeclarations(); err != nil {
		fmt.Printf("Error loading declarations: %v\n", err)
	}

	return api
}

// parseReference extracts book name, chapter, and verse from a reference string.
// For complex references (multiple verses/books), we sort using only the first
// Bible reference segment so that overall ordering is stable and predictable.
func parseReference(reference string) (book string, chapter int, verse int) {
	// Remove leading/trailing whitespace
	reference = strings.TrimSpace(reference)
	if reference == "" {
		return "", 0, 0
	}

	// Take only the first reference segment before separators like ';', 'and', '&', ', '
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
	if primary == "" {
		return reference, 0, 0
	}

	// Now parse the primary segment as "Book Chapter:Verse"
	parts := strings.Fields(primary)
	if len(parts) < 2 {
		// No obvious chapter/verse portion; treat whole thing as a non-Bible "book"
		return primary, 0, 0
	}

	// Last part should be chapter:verse
	chapterVerse := parts[len(parts)-1]
	book = strings.Join(parts[:len(parts)-1], " ")

	parseChapterVerse(chapterVerse, &chapter, &verse)
	return book, chapter, verse
}

// parseChapterVerse parses "28:8" format into chapter and verse
func parseChapterVerse(chapterVerse string, chapter *int, verse *int) {
	parts := strings.Split(chapterVerse, ":")
	if len(parts) >= 2 {
		if c, err := strconv.Atoi(parts[0]); err == nil {
			*chapter = c
		}
		if v, err := strconv.Atoi(parts[1]); err == nil {
			*verse = v
		}
	}
}

// sortDeclarationsByReference sorts declarations by Bible book order, then chapter, then verse
func (api *DeclarationsAPI) sortDeclarationsByReference() {
	sort.Slice(api.declarations, func(i, j int) bool {
		bookI, chapterI, verseI := parseReference(api.declarations[i].Reference)
		bookJ, chapterJ, verseJ := parseReference(api.declarations[j].Reference)

		orderI, existsI := bibleBookOrder[bookI]
		orderJ, existsJ := bibleBookOrder[bookJ]

		// If book order is found for both, compare by book order
		if existsI && existsJ {
			if orderI != orderJ {
				return orderI < orderJ
			}
			// Same book, compare by chapter
			if chapterI != chapterJ {
				return chapterI < chapterJ
			}
			// Same chapter, compare by verse
			return verseI < verseJ
		}

		// If only one book order is found, prioritize the known (Bible) one so
		// name/non-Bible references naturally sort to the end of the file
		if existsI && !existsJ {
			return true
		}
		if !existsI && existsJ {
			return false
		}

		// If neither book order is found, sort alphabetically by book name
		if bookI != bookJ {
			return bookI < bookJ
		}
		// Same book, compare by chapter
		if chapterI != chapterJ {
			return chapterI < chapterJ
		}
		// Same chapter, compare by verse
		return verseI < verseJ
	})
}

// parseDeclaration parses a line from the declarations file
func (api *DeclarationsAPI) parseDeclaration(line string, id int) Declaration {
	decl := Declaration{
		ID:      id,
		RawLine: line,
	}

	// Check for labels (text between colons at the start, can be multiple)
	labelRegex := regexp.MustCompile(`^:([^:]+(?::[^:]+)*?):\s*(.*)$`)
	if matches := labelRegex.FindStringSubmatch(line); matches != nil {
		decl.Label = matches[1]
		line = matches[2]
	}

	// Split declaration and reference using the " - " delimiter from the design
	parts := strings.SplitN(line, " - ", 2)
	if len(parts) == 2 {
		decl.Text = strings.TrimSpace(parts[0])
		decl.Reference = strings.TrimSpace(parts[1])
	} else {
		// If no delimiter is found, treat the whole line as text
		decl.Text = strings.TrimSpace(line)
	}

	return decl
}

// loadDeclarations loads declarations from the file
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
		// Skip empty lines and comment lines starting with '#'
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

// SortAndSaveWithComments is used in the standalone sort mode. It preserves
// comment lines (starting with '#') while re-sorting and rewriting
// declarations based on their references.
func (api *DeclarationsAPI) SortAndSaveWithComments() error {
	// Read original file to capture comment lines
	file, err := os.Open(api.filename)
	if err != nil {
		return fmt.Errorf("failed to open file for sorting: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	comments := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			comments = append(comments, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file for sorting: %w", err)
	}

	// Now sort declarations and rewrite file with comments at the top
	api.sortDeclarationsByReference()

	out, err := os.Create(api.filename)
	if err != nil {
		return fmt.Errorf("failed to create sorted file: %w", err)
	}
	defer out.Close()

	// Write preserved comments first, in original order
	for _, c := range comments {
		if _, err := fmt.Fprintln(out, c); err != nil {
			return fmt.Errorf("failed to write comment line: %w", err)
		}
	}

	// Then write sorted declarations in the same format as saveDeclarations
	for _, decl := range api.declarations {
		var line string
		if decl.Label != "" {
			line = fmt.Sprintf(":%s: %s - %s", decl.Label, decl.Text, decl.Reference)
		} else {
			line = fmt.Sprintf("%s - %s", decl.Text, decl.Reference)
		}

		if _, err := fmt.Fprintln(out, line); err != nil {
			return fmt.Errorf("failed to write declaration line: %w", err)
		}
	}

	return nil
}

// saveDeclarations saves declarations to the file
func (api *DeclarationsAPI) saveDeclarations() error {
	// Sort declarations by Bible reference before saving
	api.sortDeclarationsByReference()

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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	api.mutex.RLock()
	defer api.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"declarations": api.declarations,
	})
}

// GetDeclaration handles GET /api/v1/declarations/{id}
func (api *DeclarationsAPI) GetDeclaration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/declarations/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid declaration ID", http.StatusBadRequest)
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	api.mutex.RLock()
	defer api.mutex.RUnlock()

	if len(api.declarations) == 0 {
		http.Error(w, "No declarations available", http.StatusNotFound)
		return
	}

	// Get random declaration
	randomIndex := rand.Intn(len(api.declarations))
	randomDecl := api.declarations[randomIndex]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(randomDecl)
}

// CreateDeclaration handles POST /api/v1/declarations
func (api *DeclarationsAPI) CreateDeclaration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Label     string `json:"label"`
		Text      string `json:"text"`
		Reference string `json:"reference"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Text == "" || req.Reference == "" {
		http.Error(w, "Text and reference are required", http.StatusBadRequest)
		return
	}

	api.mutex.Lock()
	defer api.mutex.Unlock()

	newDecl := Declaration{
		ID:        api.nextID,
		Label:     req.Label,
		Text:      req.Text,
		Reference: req.Reference,
	}

	// Create raw line for consistency
	if newDecl.Label != "" {
		newDecl.RawLine = fmt.Sprintf(":%s: %s . - %s", newDecl.Label, newDecl.Text, newDecl.Reference)
	} else {
		newDecl.RawLine = fmt.Sprintf("%s . - %s", newDecl.Text, newDecl.Reference)
	}

	api.declarations = append(api.declarations, newDecl)
	api.nextID++

	if err := api.saveDeclarations(); err != nil {
		http.Error(w, "Failed to save declaration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newDecl)
}

// UpdateDeclaration handles PUT /api/v1/declarations/{id}
func (api *DeclarationsAPI) UpdateDeclaration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/declarations/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid declaration ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Label     string `json:"label"`
		Text      string `json:"text"`
		Reference string `json:"reference"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Text == "" || req.Reference == "" {
		http.Error(w, "Text and reference are required", http.StatusBadRequest)
		return
	}

	api.mutex.Lock()
	defer api.mutex.Unlock()

	decl, index := api.findDeclarationByID(id)
	if decl == nil {
		http.Error(w, "Declaration not found", http.StatusNotFound)
		return
	}

	// Update the declaration
	api.declarations[index].Label = req.Label
	api.declarations[index].Text = req.Text
	api.declarations[index].Reference = req.Reference

	// Update raw line for consistency
	if req.Label != "" {
		api.declarations[index].RawLine = fmt.Sprintf(":%s: %s . - %s", req.Label, req.Text, req.Reference)
	} else {
		api.declarations[index].RawLine = fmt.Sprintf("%s . - %s", req.Text, req.Reference)
	}

	if err := api.saveDeclarations(); err != nil {
		http.Error(w, "Failed to save declaration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.declarations[index])
}

// DeleteDeclaration handles DELETE /api/v1/declarations/{id}
func (api *DeclarationsAPI) DeleteDeclaration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/declarations/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid declaration ID", http.StatusBadRequest)
		return
	}

	api.mutex.Lock()
	defer api.mutex.Unlock()

	_, index := api.findDeclarationByID(id)
	if index == -1 {
		http.Error(w, "Declaration not found", http.StatusNotFound)
		return
	}

	// Remove the declaration
	api.declarations = append(api.declarations[:index], api.declarations[index+1:]...)

	if err := api.saveDeclarations(); err != nil {
		http.Error(w, "Failed to save declarations", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetHealth handles GET /api/v1/health
func (api *DeclarationsAPI) GetHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	api.mutex.RLock()
	defer api.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":             "healthy",
		"declarations_count": len(api.declarations),
	})
}
