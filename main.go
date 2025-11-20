package main

import (
	"log"
	"net/http"
	"strings"
)

// setupRoutes configures HTTP routes
func setupRoutes(api *DeclarationsAPI) {
	// Handle all API requests with a single handler
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Handle /api/v1/declarations (exact match) and /api/v1/declarations/
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

		// Handle /api/v1/declarations/random
		if path == "/api/v1/declarations/random" {
			api.GetRandomDeclaration(w, r)
			return
		}

		// Handle /api/v1/health
		if path == "/api/v1/health" {
			api.GetHealth(w, r)
			return
		}

		// Handle /api/v1/declarations/{id} - must have ID after /declarations/
		if strings.HasPrefix(path, "/api/v1/declarations/") && len(path) > len("/api/v1/declarations/") && path != "/api/v1/declarations/random" {
			switch r.Method {
			case http.MethodGet:
				api.GetDeclaration(w, r)
			case http.MethodPut:
				api.UpdateDeclaration(w, r)
			case http.MethodDelete:
				api.DeleteDeclaration(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// If we get here, it's an invalid API path
		http.Error(w, "API endpoint not found", http.StatusNotFound)
	})

	// UI route
	http.HandleFunc("/", ServeUI)
}

func main() {
	// Initialize the API
	api := NewDeclarationsAPI("declarations.txt")

	// Setup routes
	setupRoutes(api)

	log.Println("Starting declarations server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
