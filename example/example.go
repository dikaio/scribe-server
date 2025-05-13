package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dikaio/scribe-server/server"
)

// This is an example of how to use the scribe-server in your scribe project
func main() {
	// Get the directory to serve from command line or use default
	dir := "./public"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	// Basic usage with default options
	fmt.Println("Basic usage example:")
	fmt.Println("-------------------")
	basicServer := server.NewServer(dir)
	
	// Advanced usage with custom configuration
	fmt.Println("\nAdvanced usage example:")
	fmt.Println("---------------------")
	config := server.Config{
		RootDir:          dir,
		AddHtmlExtension: true,
		EnableLogging:    true,
		NotFoundHandler: func(path string) []byte {
			return []byte(fmt.Sprintf("<html><body><h1>Page Not Found</h1><p>The page %s could not be found.</p></body></html>", path))
		},
	}
	_ = server.NewServerWithConfig(config) // Create server but not using it in this example
	
	// Usage in scribe project
	fmt.Println("\nScribe project integration example:")
	fmt.Println("--------------------------------")
	fmt.Println("In your scribe project, after generating your site:")
	fmt.Println(`
import (
	"github.com/dikaio/scribe-server/server"
)

func startServer(outputDir string) {
	// Create a server that doesn't redirect paths without trailing slashes
	config := server.Config{
		RootDir:          outputDir,
		AddHtmlExtension: true,
		EnableLogging:    true,
	}
	srv := server.NewServerWithConfig(config)
	
	// Start the server
	log.Printf("Starting preview server at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}`)

	// Start the server using the basic configuration
	fmt.Println("\nStarting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", basicServer))
}