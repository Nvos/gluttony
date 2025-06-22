package seeds

import "embed"

//go:embed *.sql
var Seeds embed.FS
