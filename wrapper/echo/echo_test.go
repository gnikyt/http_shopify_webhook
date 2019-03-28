package echo

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
)

// Test success.
func TestGinWrapperSuccess(t *testing.T) {
	// Set our data.
	key := "secret"
	body := "{\"key\":\"value\"}"
	hmac := "ee2012a00f1649bc35f4cfe1fa582b2ebda5cbf2ef82713d6dc2ec93d81f96fb"
	shop := "example.myshopify.com"

	// Setup the server with our data.
	rec, ran := setupServer(key, shop, hmac, body)

	if c := rec.Code; c != http.StatusOK {
		t.Errorf("expected status code %v got %v", http.StatusOK, c)
	}

	if !ran {
		t.Errorf("expected next handler to run but it did not")
	}
}

// Test failure.
func TestGinWrapperFailure(t *testing.T) {
	// Set our data.
	key := "secret"
	body := "{\"key\":\"value\"}"
	hmac := "ee2012a00f1649bc35f"
	shop := "example.myshopify.com"

	// Setup the server with our data.
	rec, ran := setupServer(key, shop, hmac, body)

	if c := rec.Code; c != http.StatusBadRequest {
		t.Errorf("expected status code %v got %v", http.StatusBadRequest, c)
	}

	if ran {
		t.Errorf("expected next handler to not run but it did")
	}
}

// Sets up the server for a few tests.
func setupServer(key string, shop string, hmac string, body string) (*httptest.ResponseRecorder, bool) {
	// Setup the recorder and request.
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/webhook/order-create", bytes.NewBufferString(body))

	// Set the headers.
	req.Header.Set("X-Shopify-Shop-Domain", shop)
	req.Header.Set("X-Shopify-Hmac-Sha256", hmac)

	// Set the handler for the request.
	ran := false
	nh := func(c echo.Context) error {
		ran = true
		c.String(http.StatusOK, "Ok")

		return nil
	}

	// Start Echo.
	e := echo.New()

	// Setup the recorder and request.
	c := e.NewContext(req, rec)
	h := WebhookVerify(key)(nh)
	h(c)

	return rec, ran
}
