package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Initialize the API with declarations file
	api := NewDeclarationsAPI("declarations.txt")

	// If SORT_ONLY=1, just resort and rewrite declarations.txt, then exit
	if os.Getenv("SORT_ONLY") == "1" {
		if err := api.SortAndSave(); err != nil {
			log.Fatalf("Failed to sort declarations: %v", err)
		}
		log.Println("Sorted declarations.txt successfully")
		return
	}

	// Setup HTTP routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			ServeUI(w, r)
			return
		}

		// Handle API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			handleAPIRoutes(api, w, r)
			return
		}

		http.NotFound(w, r)
	})

	// Start server
	log.Println("Starting declarations server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleAPIRoutes(api *DeclarationsAPI, w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// GET /api/v1/health
	if path == "/api/v1/health" {
		api.GetHealth(w, r)
		return
	}

	// GET /api/v1/declarations/random
	if path == "/api/v1/declarations/random" {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		api.GetRandomDeclaration(w, r)
		return
	}

	// GET /api/v1/declarations or POST /api/v1/declarations
	if path == "/api/v1/declarations" || path == "/api/v1/declarations/" {
		switch r.Method {
		case http.MethodGet:
			api.GetDeclarations(w, r)
		case http.MethodPost:
			api.CreateDeclaration(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// GET/PUT/DELETE /api/v1/declarations/{id}
	if strings.HasPrefix(path, "/api/v1/declarations/") {
		// Extract ID
		idStr := strings.TrimPrefix(path, "/api/v1/declarations/")
		if idStr != "" && idStr != "random" {
			switch r.Method {
			case http.MethodGet:
				api.GetDeclaration(w, r, idStr)
			case http.MethodPut:
				api.UpdateDeclaration(w, r, idStr)
			case http.MethodDelete:
				api.DeleteDeclaration(w, r, idStr)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
	}

	// GET /api/v1/bible/text?q={reference}
	if path == "/api/v1/bible/text" {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		GetBibleText(w, r)
		return
	}

	http.Error(w, "API endpoint not found", http.StatusNotFound)
}
