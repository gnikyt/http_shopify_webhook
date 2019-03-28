package http_shopify_webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
)

// Public webhook verify function wrapper.
// Can be used with any framework tapping into net/http.
// Simply pass in the secret key for the Shopify app.
// Example: `WebhookVerify("abc123", anotherHandler)`.
func WebhookVerify(key string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WebhookVerifyRequest(key, w, r)
		fn(w, r)
	}
}

// Webhook verify request from HTTP.
// Returns a usable handler.
// Pass in the secret key for the Shopify app and the next handler.`
func WebhookVerifyRequest(key string, w http.ResponseWriter, r *http.Request) {
	// HMAC from request headers and the shop.
	shmac := r.Header.Get("X-Shopify-Hmac-Sha256")
	shop := r.Header.Get("X-Shopify-Shop-Domain")

	// Use TeeReader as it won't destroy the original body.
	var buffer bytes.Buffer
	tr := io.TeeReader(r.Body, &buffer)
	bb, _ := ioutil.ReadAll(tr)

	// Verify all is ok.
	if ok := verifyRequest(key, shop, shmac, bb); !ok {
		http.Error(w, "Invalid webhook signature", http.StatusBadRequest)
	}
}

// Do the actual work.
// Take the request body, the secret key,
// Attempt to reproduce the same HMAC from the request.
func verifyRequest(key string, shop string, shmac string, bb []byte) bool {
	if shop == "" {
		// No shop provided.
		return false
	}

	// Create an hmac of the body with the secret key to compare.
	h := hmac.New(sha256.New, []byte(key))
	h.Write(bb)
	sum := hex.EncodeToString(h.Sum(nil))

	return sum == shmac
}
