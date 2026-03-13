package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitMiddleware cria um middleware de rate limiting por IP
// rate: formato "X-Y" onde X é o número de requests e Y é o período (S, M, H, D)
// Exemplos:
//   - "5-M" = 5 requests por minuto
//   - "10-S" = 10 requests por segundo
//   - "100-H" = 100 requests por hora
func RateLimitMiddleware(rate string) gin.HandlerFunc {
	// Parse do rate (ex: "5-M" = 5 requests por minuto)
	parsedRate, err := limiter.NewRateFromFormatted(rate)
	if err != nil {
		// Fallback seguro: 5 tentativas por minuto
		parsedRate = limiter.Rate{
			Period: 1 * time.Minute,
			Limit:  5,
		}
	}

	// Store em memória
	// Para produção com múltiplas instâncias, considerar Redis
	store := memory.NewStore()

	// Criar limiter
	lim := limiter.New(store, parsedRate)

	// Middleware do Gin
	middleware := mgin.NewMiddleware(lim)

	return func(c *gin.Context) {
		// Aplicar o middleware do limiter
		middleware(c)

		// Verificar se o limite foi atingido
		// O middleware do limiter já seta os headers X-RateLimit-*
		// e aborta com status 429 automaticamente
	}
}
