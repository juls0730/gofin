//go:build ssr
// +build ssr

package main

import (
	"embed"
	"github.com/juls0730/gofin/ui"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
)

func init() {
	initUi = func(app *fiber.App) {
		if !fiber.IsChild() {
			tmpDir, err := os.MkdirTemp("", "gofin-ssr")
			if err != nil {
				panic(err)
			}

			err = copyEmbeddedFiles(ui.DistDir, ".output", tmpDir)
			if err != nil {
				panic(err)
			}

			path := filepath.Join(tmpDir, "server/index.mjs")

			// bun is nice, but it still has issues with memory leaks
			spawnProcess("node", []string{path}, app)
		}

		target := "localhost:3000"
		app.All("/*", func(c fiber.Ctx) error {
			path := c.Path()
			if strings.HasPrefix(path, "/api") {
				return c.Next()
			}

			request := c.Request().URI()
			return proxy.Do(c, "http://"+target+string(request.RequestURI()))
		})
	}
}

func copyEmbeddedFiles(fs embed.FS, sourcePath string, targetDir string) error {
	entries, err := fs.ReadDir(sourcePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourceFile := filepath.Join(sourcePath, entry.Name())
		destFile := filepath.Join(targetDir, entry.Name())

		if entry.IsDir() {
			os.Mkdir(destFile, 0755)
			err := copyEmbeddedFiles(fs, sourceFile, destFile)
			if err != nil {
				return err
			}
		} else {
			source, err := fs.Open(sourceFile)
			if err != nil {
				return err
			}
			defer source.Close()

			dest, err := os.Create(destFile)
			if err != nil {
				return err
			}
			defer dest.Close()

			_, err = io.Copy(dest, source)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
