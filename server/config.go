package server

// Config holds server configuration options
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

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		RootDir:          "./public",
		AddHtmlExtension: true,
		EnableLogging:    true,
		NotFoundHandler:  nil,
	}
}