package assets

import "embed"

//go:embed email/*.html
var EmailFS embed.FS
