package http_shopify_webhook

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test base verification function works.
func TestVerifyRequest(t *testing.T) {
	// Setup a simple body with a matching HMAC.
	body := []byte("{\"key\":\"value\"}")
	hmac := "ee2012a00f1649bc35f4cfe1fa582b2ebda5cbf2ef82713d6dc2ec93d81f96fb"
	shop := "example.myshopify.com"

	if ok := verifyRequest("secret", shop, hmac, body); !ok {
		t.Errorf("expected request data to verify")
	}
}

func TestVerifyRequestError(t *testing.T) {
	// Setup a simple body with a matching HMAC, but missing shop.
	body := []byte("{\"key\":\"value\"}")
	hmac := "ee2012a00f1649bc35f4cfe1fa582b2ebda5cbf2ef82713d6dc2ec93d81f96fb"
	shop := ""

	if ok := verifyRequest("secret", shop, hmac, body); ok {
		t.Errorf("expected request data to not verify, but it did")
	}

	// Now add the shop, but make the HMAC not match.
	shop = "example.myshopify.com"
	hmac = "ee2012a00f1649bc35f4"

	if ok := verifyRequest("secret", shop, hmac, body); ok {
		t.Errorf("expected request data to not verify, but it did")
	}
}

// Test the implementation with a server.
func TestNetHttpSuccess(t *testing.T) {
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
		t.Errorf("expected next handler to run but did not")
	}
}

// Test the implementation with a server (failure of HMAC).
func TestNetHttpFailure(t *testing.T) {
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

	if ran == true {
		t.Errorf("expected next handler to not run but it did")
	}
}

// Sets up the server for a few tests.
func setupServer(key string, shop string, hmac string, body string) (*httptest.ResponseRecorder, bool) {
	// Create a mock request to use
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/webhook/order-create", bytes.NewBufferString(body))

	// Set the headers.
	req.Header.Set("X-Shopify-Shop-Domain", shop)
	req.Header.Set("X-Shopify-Hmac-Sha256", hmac)

	// Our "next" handler.
	ran := false
	nh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Ok")
		ran = true
	})

	// Create the handler and serve with our recorder and request.
	h := WebhookVerify(key, nh)
	h.ServeHTTP(rec, req)

	return rec, ran
}
