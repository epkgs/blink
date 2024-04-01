package devtools

import (
	"embed"
	"io/fs"
)

//go:embed front_end/*
var frontEnd embed.FS

var FS, _ = fs.Sub(frontEnd, "front_end")
