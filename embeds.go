package main

import (
	_ "embed"
)

// embed the file templates
var (
	//go:embed structure/internal/shared/head/head.css
	css_shared string
	//go:embed structure/internal/shared/head/head.js
	js_shared string
	//go:embed structure/internal/pages/main/main.html
	page_tmpl string
	//go:embed structure/internal/shared/head/head.html
	head_tmpl string
	//go:embed structure/handlers.go
	handlers_go_tmpl string
	//go:embed structure/helpers.go
	helpers_go_tmpl string
	//go:embed structure/main.go
	main_go_tmpl string
	//go:embed structure/router.go
	router_tmpl string
	//go:embed structure/server.go
	server_go_tmpl string
	//go:embed structure/viewdata.go
	viewdata_go_tmpl string
	//go:embed structure/globals.go
	globals_go_tmpl string
	//go:embed streamable/tmpls/submitform.html
	globals_streamable_html string
	//go:embed streamable/tmpls/submithandler.go
	globals_streamable_go string
	//go:embed structure/gitignore
	gitignore string

	files map[string]string = map[string]string{
		head_tmpl:        "internal/shared/head/head.html",
		css_shared:       "internal/shared/head/head.css",
		js_shared:        "internal/shared/head/head.js",
		page_tmpl:        "internal/pages/main/main.html",
		handlers_go_tmpl: "handlers.go",
		helpers_go_tmpl:  "helpers.go",
		main_go_tmpl:     "main.go",
		router_tmpl:      "router.go",
		server_go_tmpl:   "server.go",
		viewdata_go_tmpl: "viewdata.go",
		globals_go_tmpl:  "globals.go",
		gitignore:        ".gitignore",
	}
)
