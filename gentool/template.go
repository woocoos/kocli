package gentool

import (
	"embed"
)

var (
	//go:embed template/*
	templateDir embed.FS
)
