//go:build dev
// +build dev

package main

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
	fastProxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
)

var (
	proxyServer, _ = fastProxy.NewWSReverseProxyWith(fastProxy.WithURL_OptionWS("ws://localhost:3000/_nuxt/"))
)

func init() {
	initUi = func(app *fiber.App) {
		if !fiber.IsChild() {
			spawnProcess("bun", []string{"--cwd=ui", "run", "dev"}, app)
		}

		target := "localhost:3000"
		app.Get("/_nuxt/", func(c fiber.Ctx) error {
			proxyServer.ServeHTTP(c.Context())
			return nil
		})

		app.All("/*", func(c fiber.Ctx) error {
			path := c.Path()
			if strings.HasPrefix(path, "/api") {
				return c.Next()
			}

			requestUri := string(c.Request().URI().RequestURI())

			return proxy.Do(c, "http://"+target+requestUri)
		})
	}
}
