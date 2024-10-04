package main

import (
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/pprof"
)

var initUi func(app *fiber.App)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder:                  sonic.Marshal,
		JSONDecoder:                  sonic.Unmarshal,
		StreamRequestBody:            true,
	})

	app.Use(cors.New())
	app.Use(helmet.New())

	// This is your job
	// api := app.Group("/api")
	// {
	//
	// }

	// calls out to a function set by either server.go server_dev.go based on the presence of the dev tag, and hosts
	// either the static files that get embedded into the binary in ui/embed.go or proxies the dev server that gets
	// run in the provided function
	initUi(app)

	if err := app.Listen(":1435"); err != nil && err != http.ErrServerClosed {
		fmt.Println("Error starting HTTP server:", err)
	}
}
