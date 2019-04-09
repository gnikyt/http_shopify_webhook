package http_shopify_webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

// Public webhook verify function wrapper.
// Can be used with any framework tapping into net/http.
// Simply pass in the secret key for the Shopify app.
// Example: `WebhookVerify("abc123", anotherHandler)`.
func WebhookVerify(key string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify and if all is well, run the next handler.
		ok := WebhookVerifyRequest(key, w, r)
		if ok {
			fn(w, r)
		}
	}
}

// Webhook verify request from HTTP.
// Returns a usable handler.
// Pass in the secret key for the Shopify app and the next handler.`
func WebhookVerifyRequest(key string, w http.ResponseWriter, r *http.Request) (ok bool) {
	// HMAC from request headers and the shop.
	shmac := r.Header.Get("X-Shopify-Hmac-Sha256")
	shop := r.Header.Get("X-Shopify-Shop-Domain")

	// Read the body and put it back.
	bb, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bb))

	// Verify all is ok.
	ok = verifyRequest(key, shop, shmac, bb)
	if !ok {
		http.Error(w, "Invalid webhook signature", http.StatusBadRequest)
		return
	}

	return
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
	enc := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return enc == shmac
}
