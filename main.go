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

	// GET /api/v1/env
	if path == "/api/v1/env" {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		GetEnv(w, r)
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

	// GET /api/v1/declarations (read-only)
	if path == "/api/v1/declarations" || path == "/api/v1/declarations/" {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		api.GetDeclarations(w, r)
		return
	}

	// GET /api/v1/labels
	if path == "/api/v1/labels" || path == "/api/v1/labels/" {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		api.GetLabels(w, r)
		return
	}

	// GET /api/v1/declarations/label/{label}
	if strings.HasPrefix(path, "/api/v1/declarations/label/") {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		label := strings.TrimPrefix(path, "/api/v1/declarations/label/")
		if label != "" {
			api.GetDeclarationsByLabel(w, r, label)
			return
		}
	}

	// GET /api/v1/declarations/{id}
	if strings.HasPrefix(path, "/api/v1/declarations/") {
		// Extract ID
		idStr := strings.TrimPrefix(path, "/api/v1/declarations/")
		// Skip if it's "random" or "label" - those are handled above
		if idStr != "" && idStr != "random" && !strings.HasPrefix(idStr, "label/") {
			if r.Method != http.MethodGet {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			api.GetDeclaration(w, r, idStr)
			return
		}
	}

	// GET /api/v1/bible-esv/{reference}
	if strings.HasPrefix(path, "/api/v1/bible-esv/") {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		reference := strings.TrimPrefix(path, "/api/v1/bible-esv/")
		if reference != "" {
			GetBibleText(w, r, reference)
			return
		}
	}

	http.Error(w, "API endpoint not found", http.StatusNotFound)
}
