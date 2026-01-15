package main

/*
AI Gateway - High-Performance LLM API Router
Built by Revy (˃ᆺ˂) 💜

A blazing-fast API gateway written in Go that routes requests to multiple
LLM providers with load balancing, caching, and rate limiting.
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ModelProvider represents different LLM providers
type ModelProvider string

const (
	OpenAI    ModelProvider = "openai"
	Anthropic ModelProvider = "anthropic"
	Google    ModelProvider = "google"
	DeepSeek  ModelProvider = "deepseek"
)

// LLMRequest represents an incoming request
type LLMRequest struct {
	Prompt   string        `json:"prompt"`
	Model    string        `json:"model"`
	Provider ModelProvider `json:"provider"`
	MaxTokens int          `json:"max_tokens,omitempty"`
	Temperature float64    `json:"temperature,omitempty"`
}

// LLMResponse represents the API response
type LLMResponse struct {
	Provider     ModelProvider `json:"provider"`
	Model        string        `json:"model"`
	Response     string        `json:"response"`
	TokensUsed   int           `json:"tokens_used"`
	ResponseTime float64       `json:"response_time_ms"`
	Cached       bool          `json:"cached"`
}

// Cache struct for response caching
type Cache struct {
	mu    sync.RWMutex
	data  map[string]CacheEntry
	maxSize int
}

type CacheEntry struct {
	Response  LLMResponse
	Timestamp time.Time
	TTL       time.Duration
}

// NewCache creates a new cache instance
func NewCache(maxSize int) *Cache {
	return &Cache{
		data:    make(map[string]CacheEntry),
		maxSize: maxSize,
	}
}

// Get retrieves from cache
func (c *Cache) Get(key string) (LLMResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	entry, exists := c.data[key]
	if !exists {
		return LLMResponse{}, false
	}
	
	// Check if expired
	if time.Since(entry.Timestamp) > entry.TTL {
		return LLMResponse{}, false
	}
	
	return entry.Response, true
}

// Set stores in cache
func (c *Cache) Set(key string, response LLMResponse, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Simple eviction if cache is full
	if len(c.data) >= c.maxSize {
		// Remove oldest entry
		var oldestKey string
		oldestTime := time.Now()
		for k, v := range c.data {
			if v.Timestamp.Before(oldestTime) {
				oldestTime = v.Timestamp
				oldestKey = k
			}
		}
		delete(c.data, oldestKey)
	}
	
	c.data[key] = CacheEntry{
		Response:  response,
		Timestamp: time.Now(),
		TTL:       ttl,
	}
}

// RateLimiter for API rate limiting
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow checks if request is allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	windowStart := now.Add(-rl.window)
	
	// Clean old requests
	requests := rl.requests[key]
	validRequests := []time.Time{}
	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}
	
	// Check limit
	if len(validRequests) >= rl.limit {
		return false
	}
	
	// Add new request
	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests
	
	return true
}

// Gateway is the main API gateway
type Gateway struct {
	cache       *Cache
	rateLimiter *RateLimiter
	metrics     *Metrics
}

// Metrics tracks API usage
type Metrics struct {
	mu            sync.RWMutex
	totalRequests int64
	cacheHits     int64
	cacheMisses   int64
	errors        int64
}

// NewGateway creates a new gateway instance
func NewGateway() *Gateway {
	return &Gateway{
		cache:       NewCache(1000),
		rateLimiter: NewRateLimiter(100, time.Minute),
		metrics:     &Metrics{},
	}
}

// HandleLLMRequest processes incoming LLM requests
func (g *Gateway) HandleLLMRequest(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	
	// Rate limiting
	clientIP := r.RemoteAddr
	if !g.rateLimiter.Allow(clientIP) {
		http.Error(w, `{"error":"Rate limit exceeded"}`, http.StatusTooManyRequests)
		g.metrics.RecordError()
		return
	}
	
	// Parse request
	var req LLMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		g.metrics.RecordError()
		return
	}
	
	// Generate cache key
	cacheKey := fmt.Sprintf("%s:%s:%s", req.Provider, req.Model, req.Prompt)
	
	// Check cache
	if cached, found := g.cache.Get(cacheKey); found {
		g.metrics.RecordCacheHit()
		cached.Cached = true
		json.NewEncoder(w).Encode(cached)
		return
	}
	
	g.metrics.RecordCacheMiss()
	
	// Process request
	startTime := time.Now()
	response, err := g.processLLMRequest(r.Context(), req)
	responseTime := time.Since(startTime).Milliseconds()
	
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		g.metrics.RecordError()
		return
	}
	
	response.ResponseTime = float64(responseTime)
	
	// Cache response
	g.cache.Set(cacheKey, response, 1*time.Hour)
	
	// Send response
	g.metrics.RecordRequest()
	json.NewEncoder(w).Encode(response)
}

// processLLMRequest handles the actual LLM API call
func (g *Gateway) processLLMRequest(ctx context.Context, req LLMRequest) (LLMResponse, error) {
	// This would call the actual LLM APIs
	// For demo purposes, return simulated response
	
	switch req.Provider {
	case OpenAI:
		return g.callOpenAI(ctx, req)
	case Anthropic:
		return g.callAnthropic(ctx, req)
	case Google:
		return g.callGoogle(ctx, req)
	case DeepSeek:
		return g.callDeepSeek(ctx, req)
	default:
		return LLMResponse{}, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

// Provider-specific methods (simulated for demo)
func (g *Gateway) callOpenAI(ctx context.Context, req LLMRequest) (LLMResponse, error) {
	// Simulate API call
	time.Sleep(500 * time.Millisecond)
	
	return LLMResponse{
		Provider:   OpenAI,
		Model:      req.Model,
		Response:   fmt.Sprintf("OpenAI response to: %s", req.Prompt),
		TokensUsed: len(req.Prompt) / 4,
		Cached:     false,
	}, nil
}

func (g *Gateway) callAnthropic(ctx context.Context, req LLMRequest) (LLMResponse, error) {
	time.Sleep(450 * time.Millisecond)
	
	return LLMResponse{
		Provider:   Anthropic,
		Model:      req.Model,
		Response:   fmt.Sprintf("Anthropic response to: %s", req.Prompt),
		TokensUsed: len(req.Prompt) / 4,
		Cached:     false,
	}, nil
}

func (g *Gateway) callGoogle(ctx context.Context, req LLMRequest) (LLMResponse, error) {
	time.Sleep(400 * time.Millisecond)
	
	return LLMResponse{
		Provider:   Google,
		Model:      req.Model,
		Response:   fmt.Sprintf("Google response to: %s", req.Prompt),
		TokensUsed: len(req.Prompt) / 4,
		Cached:     false,
	}, nil
}

func (g *Gateway) callDeepSeek(ctx context.Context, req LLMRequest) (LLMResponse, error) {
	time.Sleep(350 * time.Millisecond)
	
	return LLMResponse{
		Provider:   DeepSeek,
		Model:      req.Model,
		Response:   fmt.Sprintf("DeepSeek response to: %s", req.Prompt),
		TokensUsed: len(req.Prompt) / 4,
		Cached:     false,
	}, nil
}

// Metrics methods
func (m *Metrics) RecordRequest() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalRequests++
}

func (m *Metrics) RecordCacheHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheHits++
}

func (m *Metrics) RecordCacheMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheMisses++
}

func (m *Metrics) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors++
}

// HandleMetrics returns gateway metrics
func (g *Gateway) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	g.metrics.mu.RLock()
	defer g.metrics.mu.RUnlock()
	
	w.Header().Set("Content-Type", "application/json")
	
	cacheHitRate := 0.0
	total := g.metrics.cacheHits + g.metrics.cacheMisses
	if total > 0 {
		cacheHitRate = float64(g.metrics.cacheHits) / float64(total) * 100
	}
	
	metrics := map[string]interface{}{
		"total_requests": g.metrics.totalRequests,
		"cache_hits":     g.metrics.cacheHits,
		"cache_misses":   g.metrics.cacheMisses,
		"cache_hit_rate": fmt.Sprintf("%.2f%%", cacheHitRate),
		"errors":         g.metrics.errors,
	}
	
	json.NewEncoder(w).Encode(metrics)
}

// HandleHealth returns health status
func (g *Gateway) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func main() {
	gateway := NewGateway()
	
	// Setup routes
	http.HandleFunc("/api/llm", gateway.HandleLLMRequest)
	http.HandleFunc("/api/metrics", gateway.HandleMetrics)
	http.HandleFunc("/health", gateway.HandleHealth)
	
	// Static file serving for frontend
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	
	port := ":8080"
	
	fmt.Printf(`
╔═══════════════════════════════════════════════════════╗
║   🔥 AI Gateway - High-Performance LLM Router 🔥     ║
║                Built by Revy (˃ᆺ˂) 💜               ║
╠═══════════════════════════════════════════════════════╣
║  Server started on http://localhost%s              ║
║                                                       ║
║  Endpoints:                                           ║
║    POST   /api/llm     - LLM requests                ║
║    GET    /api/metrics - Gateway metrics             ║
║    GET    /health      - Health check                ║
╚═══════════════════════════════════════════════════════╝
`, port)
	
	log.Fatal(http.ListenAndServe(port, nil))
}
