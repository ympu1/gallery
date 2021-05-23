package main

import (
	"fmt"
	"os"
	//"io"
	"html/template"
	"net/http"
)

type handler struct {
	root string
	config config
}

func (h *handler) main(w http.ResponseWriter, r *http.Request) {
	filePath := h.config.RootPath + r.URL.Path
	filePathLastSym := filePath[len(filePath)-1:]
	if filePathLastSym == "/" {
		filePath = filePath[0:len(filePath)-1]
	}

	fileOrDir, err := os.Stat(filePath)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprint(w, 404)
		return
	}

	if fileOrDir.IsDir() {
		h.showFilesList(filePath, w, r)
	} else {
		http.ServeFile(w, r, filePath)
	}
}

func (h *handler) showFilesList(filePath string, w http.ResponseWriter, r *http.Request) {
	host := r.Host
	var scheme string

	if r.TLS == nil {
		scheme = "http://"
	} else {
		scheme = "https://"
	}

	data := make(map[string]interface{})
	var filesSlice []webFile
	osFiles, _ := os.ReadDir(filePath + "/.")
	for _, osFile := range osFiles {
		var file webFile
		file.fillFields(osFile, filePath, h.config)
		filesSlice = append(filesSlice, file)
	}

	data["files"] = filesSlice
	data["parentDir"] = template.URL(scheme + host + r.URL.Path)
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, data)
}

func (h *handler) static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/" + r.URL.Path)
}

func (h *handler) thumb(w http.ResponseWriter, r *http.Request) {
	filePath := getGetParam(r, "file")
	fileType := getGetParam(r, "type")

	if len(filePath) == 0 || len(fileType) == 0 {
		w.WriteHeader(404)
		fmt.Fprint(w, 404)
		return
	}

	var file webFile
	file.Path = filePath
	file.Type = fileType

	w.Write(file.getThumb(h.config))
}

func getGetParam(r *http.Request, param string) string {
	keys, ok := r.URL.Query()[param]
	if !ok || len(keys[0]) < 1 {
		return ""
	}
	return keys[0]
}