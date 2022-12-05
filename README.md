# http_shopify_webhook

[![Build Status](https://secure.travis-ci.org/gnikyt/http_shopify_webhook.png?branch=master)](http://travis-ci.org/gnikyt/http_shopify_webhook)
[![Coverage Status](https://coveralls.io/repos/github/gnikyt/http_shopify_webhook/badge.svg?branch=master)](https://coveralls.io/github/gnikyt/http_shopify_webhook?branch=master)

A middleware for validating incoming Shopify webhooks.

It can be used with any framework which speaks to `http.http.ResponseWriter`, `http.Request`, and `http.HandlerFunc`. Can be used with Go's `net/http`, `echo`, `gin, and others.

## Usage

This package provides ability to grab the shop's domain, the request HMAC, and POST body of a webhook request. Then, it will reproduce the HMAC locally with the POST body and the app's secret key to see if the data matches. Essentially (in pseudo code): `base64(hmac("secret", body)) == req_hmac`.

For more information [see Shopify's article](https://help.shopify.com/en/api/getting-started/webhooks).

### net/http

```go
package main

import (
  "fmt"
  "log"
  "net/http"

  hsw "github.com/gnikyt/http_shopify_webhook"
)

// Handler. Handle your webhook here.
func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Ok")
}

func main() {
  secret := "key" // Your secret key for the app.
  http.HandleFunc("/webhook/order-create", hsw.WebhookVerify(key, handler))
  log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Echo

```go
package main

import (
  "net/http"

  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  hsw "github.com/gnikyt/http_shopify_webhook/wrapper/echo"
)

// Handler. Handle your webhook here.
func hello(c echo.Context) error {
  return c.String(http.StatusOK, "Ok")
}

func main() {
  secret := "key" // Your secret key for the app.

  e := echo.New()
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.Use(hsw.WebhookVerify(secret))

  e.POST("/webhook/order-create", handler)

  e.Logger.Fatal(e.Start(":1323"))
}
```

### Gin

```go
package main

import (
  "github.com/gin-gonic/gin"

  hsw "github.com/gnikyt/http_shopify_webhook/wrapper/gin"
)

func main() {
  secret := "key" // Your secret key for the app.

  r := gin.Default()
  r.Use(hsw.WebhookVerify(secret))

  r.POST("/webhook/order-create", func(c *gin.Context) {
    // Handle your webhook here.
    c.Data(http.StatusOK, "text/plain", "Ok")
  })

  r.Run()
}
```

## Testing

`go test ./...`, fully tested.

## Documentation

Available through [godoc.org](https://godoc.org/github.com/gnikyt/http_shopify_webhook).

## LICENSE

This project is released under the MIT [license](https://github.com/gnikyt/http_shopify_webhook/blob/master/LICENSE).
