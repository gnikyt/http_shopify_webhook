package http_shopify_webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

// HTTP errors.
const (
	errMissingSignature = "Missing signature from request"
	errMissingShop      = "Missing shop from request"
	errSignature        = "Invalid webhook signature"
)

// HTTP headers.
const (
	headerHmac = "X-Shopify-Hmac-Sha256"
	headerShop = "X-Shopify-Shop-Domain"
)

// Public webhook verify wrapper.
// Can be used with any framework tapping into net/http.
// Simply pass in the secret key for the Shopify app.
// Example: `WebhookVerify("abc123", anotherHandler)`.
func WebhookVerify(key string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify and if all is well, run the next handler.
		if ok := WebhookVerifyRequest(key, w, r); ok {
			fn(w, r)
		}
	}
}

// Webhook verify request from HTTP.
// Returns a usable handler.
// Pass in the secret key for the Shopify app and the next handler.`
func WebhookVerifyRequest(key string, w http.ResponseWriter, r *http.Request) bool {
	// HMAC from request headers and the shop.
	shmac := r.Header.Get(headerHmac)
	shop := r.Header.Get(headerShop)

	if shop == "" {
		// No shop provided.
		httpError(errMissingShop, http.StatusBadRequest, w)
		return false
	}

	if shmac == "" {
		// No HMAC provided.
		httpError(errMissingSignature, http.StatusBadRequest, w)
		return false
	}

	// Read the body and put it back.
	bb, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bb))

	// Create a signature to compare.
	lhmac := newSignature(key, bb)

	// Verify all is ok.
	if ok := isValidSignature(lhmac, shmac); !ok {
		httpError(errSignature, http.StatusBadRequest, w)
		return false
	}
	return true
}

// Create an HMAC of the body (bb) with the secret key (key).
// Returns a string.
func newSignature(key string, bb []byte) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(bb)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Compares the created HMAC signature with the request's HMAC signature.
// Returns bool of comparison result.
func isValidSignature(lhmac string, shmac string) bool {
	return lhmac == shmac
}

// Common place to issue an error to the response writer.
func httpError(msg string, code int, w http.ResponseWriter) {
	http.Error(w, msg, code)
}
