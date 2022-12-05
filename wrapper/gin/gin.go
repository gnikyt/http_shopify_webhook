package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gnikyt/http_shopify_webhook"
)

// Compatible wrapper for Gin framework.
// Example: `g.Use(WebhookVerify("secret")`.
func WebhookVerify(key string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ok := http_shopify_webhook.WebhookVerifyRequest(key, ctx.Writer, ctx.Request)
		if !ok {
			ctx.AbortWithStatus(http.StatusBadRequest)
		}
	}
}
