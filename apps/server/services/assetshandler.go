package services

import (
	"embed"
	"fmt"
	"net/http"
	"path/filepath"
)

func RegisterStaticRouter(router map[string]http.HandlerFunc, fs []embed.FS) {
	if len(fs) < 2 {
		fmt.Println("RegisterStaticRouter not registered")
		return
	}
	router[""] = handleIndex("resources/index.html", fs[0])
	router["/"] = handleIndex("resources/index.html", fs[0])
	router["/index.html"] = handleIndex("resources/index.html", fs[0])
	router["/favicon.ico"] = handleIndex("resources/favicon.ico", fs[0])
	files, err := fs[1].ReadDir("resources/assets")
	if err != nil {
		fmt.Println("load static assets error:", err.Error())
		return
	}
	for _, f := range files {
		router["/assets/"+f.Name()] = handleAssets("resources/assets/"+f.Name(), fs[1])
	}
}

func handleIndex(fileName string, fs embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexHTML, err := fs.ReadFile(fileName)
		if err != nil {
			// http.Error(w, "404 Not Found", http.StatusNotFound)
			fmt.Println(err.Error())
			return
		}
		switch filepath.Ext(fileName) {
		case ".css":
			w.Header().Add("Content-Type", "text/css")
		case ".js":
			w.Header().Add("Content-Type", "text/javascript")
		case ".html":
			w.Header().Add("Content-Type", "text/html")
		}
		w.Write(indexHTML)
	}
}

func handleAssets(fileName string, fs embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexHTML, err := fs.ReadFile(fileName)
		if err != nil {
			// http.Error(w, "404 Not Found", http.StatusNotFound)
			fmt.Println(err.Error())
			return
		}
		switch filepath.Ext(fileName) {
		case ".css":
			w.Header().Add("Content-Type", "text/css")
		case ".js":
			w.Header().Add("Content-Type", "text/javascript")
		case ".html":
			w.Header().Add("Content-Type", "text/html")
		}
		w.Write(indexHTML)
	}
}
