package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alireza-akbarzadeh/ginflow/cmd/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRouter creates a router with embedded templates for testing
func setupTestRouter() *gin.Engine {
	router := gin.New()
	router.SetHTMLTemplate(template.Must(template.ParseFS(web.Templates, "pages/*.html")))
	return router
}

// ============================================================================
// Web Handler Tests
// ============================================================================

func TestShowLandingPage(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "successful landing page render",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &Handler{}

			router := setupTestRouter()
			router.GET("/", handler.ShowLandingPage)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
		})
	}
}

func TestShowHealthPage(t *testing.T) {
	tests := []struct {
		name           string
		acceptHeader   string
		expectedStatus int
		expectedType   string
	}{
		{
			name:           "HTML health page",
			acceptHeader:   "text/html",
			expectedStatus: http.StatusOK,
			expectedType:   "text/html",
		},
		{
			name:           "JSON health response",
			acceptHeader:   "application/json",
			expectedStatus: http.StatusOK,
			expectedType:   "application/json",
		},
		{
			name:           "default to HTML",
			acceptHeader:   "",
			expectedStatus: http.StatusOK,
			expectedType:   "text/html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &Handler{}

			router := setupTestRouter()
			router.GET("/health", handler.ShowHealthPage)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			if tt.acceptHeader != "" {
				req.Header.Set("Accept", tt.acceptHeader)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Header().Get("Content-Type"), tt.expectedType)

			if tt.acceptHeader == "application/json" {
				assert.Contains(t, w.Body.String(), "status")
				assert.Contains(t, w.Body.String(), "ok")
			}
		})
	}
}

func TestShowDashboardPage(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "successful dashboard page render",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &Handler{}

			router := setupTestRouter()
			router.GET("/dashboard", handler.ShowDashboardPage)

			req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
		})
	}
}

// ============================================================================
// Content Type Negotiation Tests
// ============================================================================

func TestContentTypeNegotiation(t *testing.T) {
	tests := []struct {
		name         string
		acceptHeader string
		endpoint     string
		validate     func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:         "JSON request to health endpoint",
			acceptHeader: "application/json",
			endpoint:     "/health",
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
				body := w.Body.String()
				assert.Contains(t, body, `"status"`)
				assert.Contains(t, body, `"ok"`)
			},
		},
		{
			name:         "HTML request to health endpoint",
			acceptHeader: "text/html",
			endpoint:     "/health",
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
			},
		},
		{
			name:         "browser request (multiple accept types)",
			acceptHeader: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			endpoint:     "/health",
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &Handler{}

			router := setupTestRouter()
			router.GET("/health", handler.ShowHealthPage)

			req := httptest.NewRequest(http.MethodGet, tt.endpoint, nil)
			req.Header.Set("Accept", tt.acceptHeader)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			tt.validate(t, w)
		})
	}
}

// ============================================================================
// Integration Tests for Web Pages
// ============================================================================

func TestWebPagesIntegration(t *testing.T) {
	handler := &Handler{}

	router := setupTestRouter()

	// Register all web routes
	router.GET("/", handler.ShowLandingPage)
	router.GET("/health", handler.ShowHealthPage)
	router.GET("/dashboard", handler.ShowDashboardPage)

	tests := []struct {
		name     string
		endpoint string
		method   string
		status   int
	}{
		{"landing page", "/", http.MethodGet, http.StatusOK},
		{"health page", "/health", http.MethodGet, http.StatusOK},
		{"dashboard page", "/dashboard", http.MethodGet, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.endpoint, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.status, w.Code)
			assert.NotEmpty(t, w.Body.String())
		})
	}
}

// ============================================================================
// HTTP Method Tests
// ============================================================================

func TestHTTPMethods(t *testing.T) {
	handler := &Handler{}

	router := setupTestRouter()
	router.GET("/", handler.ShowLandingPage)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{"GET allowed", http.MethodGet, http.StatusOK},
		{"POST not allowed", http.MethodPost, http.StatusNotFound},
		{"PUT not allowed", http.MethodPut, http.StatusNotFound},
		{"DELETE not allowed", http.MethodDelete, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// ============================================================================
// Response Header Tests
// ============================================================================

func TestResponseHeaders(t *testing.T) {
	handler := &Handler{}

	router := setupTestRouter()
	router.GET("/", handler.ShowLandingPage)
	router.GET("/health", handler.ShowHealthPage)

	t.Run("HTML content-type for landing page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		contentType := w.Header().Get("Content-Type")
		assert.Contains(t, contentType, "text/html")
	})

	t.Run("JSON content-type for health API", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		contentType := w.Header().Get("Content-Type")
		assert.Contains(t, contentType, "application/json")
	})
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestPageNotFound(t *testing.T) {
	router := gin.New()

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkShowLandingPage(b *testing.B) {
	handler := &Handler{}

	router := setupTestRouter()
	router.GET("/", handler.ShowLandingPage)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkShowHealthPageJSON(b *testing.B) {
	handler := &Handler{}

	router := gin.New()
	router.GET("/health", handler.ShowHealthPage)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// ============================================================================
// Table-Driven Tests for All Pages
// ============================================================================

func TestAllWebPages(t *testing.T) {
	handler := &Handler{}

	pages := []struct {
		name        string
		path        string
		handlerFunc gin.HandlerFunc
	}{
		{"landing", "/", handler.ShowLandingPage},
		{"health", "/health", handler.ShowHealthPage},
		{"dashboard", "/dashboard", handler.ShowDashboardPage},
	}

	for _, page := range pages {
		t.Run(page.name, func(t *testing.T) {
			router := setupTestRouter()
			router.GET(page.path, page.handlerFunc)

			req := httptest.NewRequest(http.MethodGet, page.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code, "Page %s should return 200", page.name)
			require.NotEmpty(t, w.Body.String(), "Page %s should have content", page.name)
		})
	}
}
