package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

func main() {
	// Determine the declarations filename based on environment
	declarationsFile, err := getDeclarationsFile()
	if err != nil {
		log.Fatalf("Failed to get declarations file: %v", err)
	}

	// Validate ESV API token is available before starting
	if err := validateESVToken(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	// Initialize the API with declarations file
	api := NewDeclarationsAPI(declarationsFile)

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

// isCloudRun detects if the app is running in Google Cloud Run
func isCloudRun() bool {
	// K_SERVICE is set in Cloud Run environment
	return os.Getenv("K_SERVICE") != ""
}

// getDeclarationsFile determines the declarations filename and downloads from GCS if needed
func getDeclarationsFile() (string, error) {
	if isCloudRun() {
		// Running in Cloud Run - download from GCS bucket
		bucketName := os.Getenv("DECLARATIONS_BUCKET_NAME")
		if bucketName == "" {
			return "", fmt.Errorf("DECLARATIONS_BUCKET_NAME environment variable not set")
		}

		// Download declarations file from bucket
		if err := downloadFromGCS(bucketName, "declarations", "declarations"); err != nil {
			return "", fmt.Errorf("failed to download declarations from GCS: %w", err)
		}

		log.Printf("Downloaded declarations file from GCS bucket: %s", bucketName)
		return "declarations", nil
	}

	// Running locally - use declarations.txt
	return "declarations.txt", nil
}

// downloadFromGCS downloads a file from Google Cloud Storage
func downloadFromGCS(bucketName, objectName, destFile string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	reader, err := obj.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to open GCS object: %w", err)
	}
	defer reader.Close()

	file, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("failed to copy GCS object to file: %w", err)
	}

	return nil
}

// validateESVToken validates that the ESV API token is available
func validateESVToken() error {
	// If running in Cloud Run, expect ESV_API_TOKEN environment variable to be set
	if isCloudRun() {
		if os.Getenv("ESV_API_TOKEN") == "" {
			return fmt.Errorf("ESV_API_TOKEN environment variable not set (should be populated from Secret Manager)")
		}
		return nil
	}

	// Running locally - check for .esv-api-token file
	if _, err := os.Stat(".esv-api-token"); os.IsNotExist(err) {
		return fmt.Errorf(".esv-api-token file not found - please create this file with your ESV API token")
	}

	return nil
}
