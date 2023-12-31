package sql

import (
	"embed"
)

//go:embed "migrations"
var EmbeddedFiles embed.FS
