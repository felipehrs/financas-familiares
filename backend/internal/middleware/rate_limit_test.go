package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipehrs/financas-familiares/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup: criar rota com rate limit de 3 requests por minuto
	r := gin.New()
	r.POST("/test", middleware.RateLimitMiddleware("3-M"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	t.Run("permite requests dentro do limite", func(t *testing.T) {
		// Primeiras 3 requests devem passar
		for i := 0; i < 3; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", nil)
			req.RemoteAddr = "192.168.1.1:1234" // Simular mesmo IP
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)

			// Verificar headers de rate limiting
			assert.NotEmpty(t, w.Header().Get("X-RateLimit-Limit"))
			assert.NotEmpty(t, w.Header().Get("X-RateLimit-Remaining"))
		}
	})

	t.Run("bloqueia 4ª request (excede limite)", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		r.ServeHTTP(w, req)

		// 4ª request deve ser bloqueada com HTTP 429
		assert.Equal(t, http.StatusTooManyRequests, w.Code, "4th request should be rate limited")

		// Header de remaining deve ser 0
		assert.Equal(t, "0", w.Header().Get("X-RateLimit-Remaining"))
	})

	t.Run("IPs diferentes têm limites independentes", func(t *testing.T) {
		// Novo IP deve ter seu próprio limite
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "192.168.1.2:5678" // IP diferente
		r.ServeHTTP(w, req)

		// Deve passar (primeiro request deste IP)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
