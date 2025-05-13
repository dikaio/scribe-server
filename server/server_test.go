package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestServer(t *testing.T) {
	// Create a temp directory for test files
	tempDir, err := os.MkdirTemp("", "scribe-server-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	createTestFiles(t, tempDir)

	// Create server with default settings
	srv := NewServer(tempDir)

	// Test cases
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Root path should serve index.html",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Index Page",
		},
		{
			name:           "Path without extension should serve .html file",
			path:           "/about",
			expectedStatus: http.StatusOK,
			expectedBody:   "About Page",
		},
		{
			name:           "Path with trailing slash should work",
			path:           "/about/",
			expectedStatus: http.StatusOK,
			expectedBody:   "About Page",
		},
		{
			name:           "Path with .html extension should work",
			path:           "/about.html",
			expectedStatus: http.StatusOK,
			expectedBody:   "About Page",
		},
		{
			name:           "Non-existent page should return 404",
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found",
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()

			srv.ServeHTTP(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if !contains(string(body), tc.expectedBody) {
				t.Errorf("Expected body to contain '%s', got '%s'", tc.expectedBody, string(body))
			}
		})
	}
}

func TestServerWithCustomConfig(t *testing.T) {
	// Create a temp directory for test files
	tempDir, err := os.MkdirTemp("", "scribe-server-custom-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	createTestFiles(t, tempDir)

	// Create custom 404 handler
	customNotFound := func(path string) []byte {
		return []byte("Custom 404: " + path + " not found")
	}

	// Create server with custom config
	config := Config{
		RootDir:          tempDir,
		AddHtmlExtension: false, // Disable .html extension
		EnableLogging:    false,
		NotFoundHandler:  customNotFound,
	}
	srv := NewServerWithConfig(config)

	// Test cases
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Path without extension should NOT serve .html file when disabled",
			path:           "/about",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Custom 404: /about not found",
		},
		{
			name:           "Path with .html extension should still work",
			path:           "/about.html",
			expectedStatus: http.StatusOK,
			expectedBody:   "About Page",
		},
		{
			name:           "Custom 404 handler should be used",
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Custom 404: /nonexistent not found",
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()

			srv.ServeHTTP(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if !contains(string(body), tc.expectedBody) {
				t.Errorf("Expected body to contain '%s', got '%s'", tc.expectedBody, string(body))
			}
		})
	}
}

// Helper functions
func createTestFiles(t *testing.T, dir string) {
	// Create index.html
	indexHTML := "<html><body>Index Page</body></html>"
	err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(indexHTML), 0644)
	if err != nil {
		t.Fatalf("Failed to create index.html: %v", err)
	}

	// Create about.html
	aboutHTML := "<html><body>About Page</body></html>"
	err = os.WriteFile(filepath.Join(dir, "about.html"), []byte(aboutHTML), 0644)
	if err != nil {
		t.Fatalf("Failed to create about.html: %v", err)
	}

	// Create a nested directory with index.html
	nestedDir := filepath.Join(dir, "nested")
	err = os.Mkdir(nestedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}

	// Create nested/index.html
	nestedHTML := "<html><body>Nested Page</body></html>"
	err = os.WriteFile(filepath.Join(nestedDir, "index.html"), []byte(nestedHTML), 0644)
	if err != nil {
		t.Fatalf("Failed to create nested/index.html: %v", err)
	}
}

// Contains checks if a string contains another string
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}