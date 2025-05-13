package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Server is a custom HTTP server that doesn't redirect paths without trailing slashes
type Server struct {
	config Config
	logger *log.Logger
}

// NewServer creates a new server that serves files from the given directory
func NewServer(rootDir string) *Server {
	return &Server{
		config: Config{
			RootDir:          rootDir,
			AddHtmlExtension: true,
			EnableLogging:    true,
		},
		logger: log.New(os.Stdout, "[SERVER] ", log.LstdFlags),
	}
}

// NewServerWithConfig creates a new server with the provided configuration
func NewServerWithConfig(config Config) *Server {
	var logger *log.Logger
	if config.EnableLogging {
		logger = log.New(os.Stdout, "[SERVER] ", log.LstdFlags)
	} else {
		logger = log.New(io.Discard, "", 0)
	}
	
	return &Server{
		config: config,
		logger: logger,
	}
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Clean the path to prevent directory traversal attacks
	path := filepath.Clean(r.URL.Path)
	
	// If it's the root path, serve index.html
	if path == "/" {
		path = "/index.html"
	}
	
	s.logger.Printf("Handling request for: %s", path)
	
	// Remove trailing slash if it exists (except for root path)
	// This is the key to preventing redirects to trailing slashes
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	
	// Try different file possibilities in this order:
	// 1. Exact path match
	// 2. Adding .html extension if no extension exists (if enabled)
	// 3. Look for index.html in directory (only if path originally had trailing slash)
	
	// Construct the full file path
	fullPath := filepath.Join(s.config.RootDir, path)
	
	// First try: exact path match
	if s.tryServeFile(w, r, fullPath) {
		return
	}
	
	// Second try: add .html extension if no extension exists and option is enabled
	if s.config.AddHtmlExtension && filepath.Ext(path) == "" {
		htmlPath := fullPath + ".html"
		if s.tryServeFile(w, r, htmlPath) {
			return
		}
	}
	
	// Third try: check for index.html in directory
	// Only do this for paths that originally had a trailing slash
	if strings.HasSuffix(r.URL.Path, "/") {
		indexPath := filepath.Join(fullPath, "index.html")
		if s.tryServeFile(w, r, indexPath) {
			return
		}
	}
	
	// If all attempts fail, handle 404
	s.logger.Printf("Not found: %s", path)
	
	// Use custom 404 handler if provided
	if s.config.NotFoundHandler != nil {
		content := s.config.NotFoundHandler(path)
		w.WriteHeader(http.StatusNotFound)
		w.Write(content)
		return
	}
	
	// Otherwise use standard 404
	http.NotFound(w, r)
}

// tryServeFile attempts to serve the file at the given path
// Returns true if successful, false otherwise
func (s *Server) tryServeFile(w http.ResponseWriter, r *http.Request, fullPath string) bool {
	info, err := os.Stat(fullPath)
	
	// If file doesn't exist or is a directory, return false
	if err != nil || info.IsDir() {
		return false
	}
	
	// Open the file
	file, err := os.Open(fullPath)
	if err != nil {
		s.logger.Printf("Error opening file %s: %v", fullPath, err)
		return false
	}
	defer file.Close()
	
	// Set content type
	contentType := getContentType(filepath.Ext(fullPath))
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	
	// Copy file content to response writer
	s.logger.Printf("Serving file: %s", fullPath)
	_, err = io.Copy(w, file)
	if err != nil {
		s.logger.Printf("Error serving file %s: %v", fullPath, err)
	}
	
	return true
}

// getContentType returns the MIME type based on file extension
func getContentType(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".xml":
		return "application/xml"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	default:
		return ""
	}
}