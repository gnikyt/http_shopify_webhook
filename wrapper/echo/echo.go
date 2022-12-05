package echo

import (
	"net/http"

	"github.com/gnikyt/http_shopify_webhook"
	"github.com/labstack/echo"
)

// Compatible wrapper for Echo framework.
// Example: `e.Use(WebhookVerify("secret"))`.
func WebhookVerify(key string) func(n echo.HandlerFunc) echo.HandlerFunc {
	return func(n echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ok := http_shopify_webhook.WebhookVerifyRequest(key, c.Response(), c.Request())
			if !ok {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid webhook request")
			}

			return n(c)
		}
	}
}
