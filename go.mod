module github.com/randomaigirl/ai-toolkit

go 1.21

require (
	// No external dependencies needed for the basic version!
	// The gateway uses only Go standard library
)

// Optional dependencies for production:
// github.com/go-redis/redis/v8 v8.11.5  // For distributed caching
// github.com/gorilla/mux v1.8.1         // For advanced routing
// github.com/prometheus/client_golang   // For metrics
