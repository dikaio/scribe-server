package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/dikaio/scribe-server/server"
)

func main() {
	// Command line flags
	port := flag.String("port", "8080", "Port to run the server on")
	dir := flag.String("dir", "./public", "Directory to serve files from")
	noHtmlExt := flag.Bool("no-html-ext", false, "Disable .html extension auto-adding")
	noLogging := flag.Bool("no-logging", false, "Disable request logging")
	flag.Parse()

	// Create a server config
	config := server.Config{
		RootDir:          *dir,
		AddHtmlExtension: !*noHtmlExt,
		EnableLogging:    !*noLogging,
	}

	// Create a new server with config
	srv := server.NewServerWithConfig(config)

	// Start the server
	log.Printf("Starting server on port %s serving files from %s\n", *port, *dir)
	log.Printf("HTML extension auto-adding: %v, Logging: %v\n", !*noHtmlExt, !*noLogging)
	log.Fatal(http.ListenAndServe(":"+*port, srv))
}