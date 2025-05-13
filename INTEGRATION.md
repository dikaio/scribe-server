# Scribe Server Implementation and Integration Guide

This document provides a detailed explanation of how we created the custom web server for Scribe and how to integrate it into the main Scribe project.

## Table of Contents

1. [Problem Statement](#problem-statement)
2. [Server Implementation Overview](#server-implementation-overview)
3. [Key Components](#key-components)
4. [How It Prevents Trailing Slash Redirects](#how-it-prevents-trailing-slash-redirects)
5. [Integration with Scribe Project](#integration-with-scribe-project)
6. [Testing](#testing)
7. [Future Enhancements](#future-enhancements)

## Problem Statement

The standard Go HTTP server (`http.FileServer`) automatically redirects URL paths that represent directories without trailing slashes to the same path with a trailing slash. For example, a request to `/about` would be redirected to `/about/` if `/about` is a directory.

This behavior is built into Go's standard library and cannot be easily changed when using `http.FileServer`. For Scribe, we wanted URLs without trailing slashes (e.g., `/about`) to be served directly without this redirection.

## Server Implementation Overview

We created a custom HTTP server that implements the `http.Handler` interface and handles file serving manually, giving us complete control over URL parsing and response generation. This approach allows us to serve files without the unwanted redirects.

The implementation consists of:

1. A custom `Server` struct that implements `http.Handler`
2. Configuration options for customizing behavior
3. Custom file lookup logic
4. Direct file serving without relying on Go's built-in file server

## Key Components

### Server Structure

The server is organized as follows:

- `server.go`: Main server implementation
- `config.go`: Configuration structures
- `server_test.go`: Test cases

### Configuration Options

Our server provides several configuration options:

```go
type Config struct {
    // RootDir is the directory to serve files from
    RootDir string
    
    // AddHtmlExtension determines whether to try adding .html extension to paths without extensions
    AddHtmlExtension bool
    
    // EnableLogging controls whether to log requests
    EnableLogging bool
    
    // NotFoundHandler is a custom handler for 404 errors
    NotFoundHandler func(string) []byte
}
```

### File Serving Logic

The file serving logic follows this priority:

1. First try: Exact path match
2. Second try: Add `.html` extension if no extension exists and the option is enabled
3. Third try: Check for `index.html` in the directory (only if the original path had a trailing slash)

If all attempts fail, the server returns a 404 response, either using a custom handler if provided or the standard `http.NotFound`.

## How It Prevents Trailing Slash Redirects

The key to preventing trailing slash redirects is in the `ServeHTTP` method:

```go
// Remove trailing slash if it exists (except for root path)
// This is the key to preventing redirects to trailing slashes
if len(path) > 1 && strings.HasSuffix(path, "/") {
    path = path[:len(path)-1]
}
```

This code explicitly removes trailing slashes from all paths (except the root path "/") before processing them further. This means that a request for both `/about` and `/about/` will be handled the same way, without any redirects.

Additionally, by manually opening files and streaming their contents to the response writer, we bypass any automatic redirect logic in the standard library:

```go
file, err := os.Open(fullPath)
// ...
io.Copy(w, file)
```

## Integration with Scribe Project

To integrate this server with the Scribe project, follow these steps:

### 1. Add as a Dependency

Update the Scribe project's `go.mod` file to include this server as a dependency:

```bash
cd /path/to/scribe/project
go get github.com/dikaio/scribe-server
```

### 2. Import the Server Package

In the Scribe project, import the server package:

```go
import "github.com/dikaio/scribe-server/server"
```

### 3. Replace the Existing Server

Find the code in Scribe that sets up the HTTP server. It likely uses `http.FileServer` directly. Replace it with our custom server:

**Before:**
```go
// This is likely what the Scribe project currently uses
func startServer(outputDir string) {
    fileServer := http.FileServer(http.Dir(outputDir))
    http.Handle("/", fileServer)
    log.Printf("Starting preview server at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**After:**
```go
// Replace with our custom server
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

### 4. Add Configuration Options (Optional)

You may want to expose configuration options to Scribe users. Add these to your configuration file or command line flags:

```go
// Example of adding server options to your command line flags
serveCmd := &cobra.Command{
    Use:   "serve",
    Short: "Start a local web server to preview the site",
    Run: func(cmd *cobra.Command, args []string) {
        config := server.Config{
            RootDir:          outputDir,
            AddHtmlExtension: viper.GetBool("server.addHtmlExtension"),
            EnableLogging:    viper.GetBool("server.enableLogging"),
        }
        srv := server.NewServerWithConfig(config)
        log.Printf("Starting preview server at http://localhost:%s", port)
        log.Fatal(http.ListenAndServe(":"+port, srv))
    },
}

serveCmd.Flags().Bool("add-html-ext", true, "Automatically add .html extension to URLs without extensions")
serveCmd.Flags().Bool("enable-logging", true, "Enable server request logging")
viper.BindPFlag("server.addHtmlExtension", serveCmd.Flags().Lookup("add-html-ext"))
viper.BindPFlag("server.enableLogging", serveCmd.Flags().Lookup("enable-logging"))
```

## Testing

To verify that the integration works correctly:

1. Start the Scribe server
2. Visit a URL without a trailing slash (e.g., `/about`)
3. Check that it serves the content directly without redirecting
4. Check server logs to confirm the request was handled as expected

Example test cases:

- `/about` should serve the content of `/about.html` directly
- `/about/` should also work and serve the same content
- URLs with extensions like `/style.css` should work normally
- The root URL `/` should serve `index.html`

## Future Enhancements

Potential improvements to consider:

1. **Cache Control Headers**: Add options for configuring cache headers
2. **MIME Type Expansion**: Expand the list of supported MIME types
3. **Custom Error Pages**: Add more options for custom error pages beyond 404
4. **HTTPS Support**: Add easy HTTPS configuration
5. **CORS Support**: Add Cross-Origin Resource Sharing headers
6. **Compression**: Add support for gzip/brotli compression

## Conclusion

This custom server implementation provides a solution to the trailing slash redirect issue while maintaining compatibility with Scribe's existing functionality. By replacing the standard `http.FileServer` with our custom implementation, we gain full control over URL handling and can serve content exactly as intended.

The implementation is lightweight, efficient, and highly configurable, making it a suitable replacement for Scribe's current server implementation.