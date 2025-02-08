package assets

import "embed"

//go:embed css/* font/* img/*
var Assets embed.FS
