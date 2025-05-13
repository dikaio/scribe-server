# Scribe Server

A simple web server for the [Scribe](https://github.com/dikaio/scribe) static site generator that doesn't redirect URLs without trailing slashes.

## Features

- Serves static files without redirecting paths like `/about` to `/about/`
- Automatic `.html` extension support (you can access `/about` instead of `/about.html`)
- Configurable logging
- Custom 404 handler support
- Designed to be imported and used in the Scribe static site generator

## Usage

### As a standalone server

```bash
# Run the server with default settings
go run main.go

# Run the server with custom options
go run main.go --port 3000 --dir ./site --no-html-ext --no-logging
```

### Command-line options

- `--port`: Port to run the server on (default: 8080)
- `--dir`: Directory to serve files from (default: ./public)
- `--no-html-ext`: Disable automatic .html extension
- `--no-logging`: Disable request logging

### As a library in your Scribe project

```go
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
}
```

## How It Works

The server handles requests for paths without trailing slashes in the following order:

1. Tries to serve the exact path as requested
2. If no file extension, tries adding `.html` to the path
3. Only for paths ending with `/`, tries to serve `index.html` in that directory

This approach ensures that URLs like `/about` will work correctly without redirecting to `/about/`.

## License

MIT