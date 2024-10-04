//go:build !dev && !ssr
// +build !dev,!ssr

package main

import (
	"github.com/juls0730/gofin/ui"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type embeddedFS struct {
	baseFS fs.FS
	prefix string
}

func (fs *embeddedFS) Open(name string) (fs.File, error) {
	// Prepend the prefix to the requested file name
	publicPath := filepath.Join(fs.prefix, name)
	file, err := fs.baseFS.Open(publicPath)
	if err != nil {
		return nil, fiber.ErrNotFound
	}

	return file, err
}

var publicFS = &embeddedFS{
	baseFS: ui.DistDir,
	prefix: ".output/public/",
}

func init() {
	initUi = func(app *fiber.App) {
		app.Use("/", static.New("", static.Config{
			FS:            publicFS,
			CacheDuration: 10 * time.Hour,
		}))

		// 404 handler
		app.Use(func(c fiber.Ctx) error {
			err := c.Next()
			if err == nil {
				return nil
			}

			path := c.Path()
			if !strings.HasPrefix(path, "/api") {
				file, err := publicFS.Open("404.html")
				if err != nil {
					c.App().Server().Logger.Printf("Error opening 404.html: %s", err)
					return err
				}
				defer file.Close()

				fileInfo, err := file.Stat()
				if err != nil {
					c.App().Server().Logger.Printf("An error occurred while getting the file info: %s", err)
					return err
				}

				fileBuf := make([]byte, fileInfo.Size())
				_, err = file.Read(fileBuf)
				if err != nil {
					c.App().Server().Logger.Printf("An error occurred while reading the file: %s", err)
					return err
				}

				c.Context().SetContentType("text/html")
				return c.Status(http.StatusOK).SendString(string(fileBuf))
			}
			return err
		})
	}
}
