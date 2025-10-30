package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var web embed.FS

var webFS, _ = fs.Sub(web, "dist")

var Handler = http.FileServerFS(webFS)
