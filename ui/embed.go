// Package ui handles the frontend embedding
package ui

import (
	"embed"
)

//go:embed all:.output
var DistDir embed.FS
