package web

import "embed"

//go:embed components/*.html pages/*.html
var Templates embed.FS
