package main

import (
	"os"
)

func newRoute(params []string) {
	route := params[0]
	handler := params[1]
	routeLine := "\t" + `mux.HandleFunc("` + route + `", ` + handler + `)`
	insertLineAfter("router.go", "func registerRoutes(mux *http.ServeMux)", routeLine, "}")
	newHandler(handler, nil)
}

func newHandler(handlerName string, fileBytes []byte) {
	f, err := os.OpenFile(handlerName+".go", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var txt string
	if fileBytes == nil {
		txt = "package main\n\nimport \"net/http\"\n\nfunc " + handlerName +
			"\n\n(w http.ResponseWriter, r *http.Request) {\n\n}"
	} else {
		txt = string(fileBytes)
	}

	if _, err = f.WriteString(txt); err != nil {
		panic(err)
	}
}
