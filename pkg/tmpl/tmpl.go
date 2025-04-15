package tmpl

import "embed"

//go:embed *.tmpl
var Tmpl embed.FS
