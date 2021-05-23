package main

import (
	"net/http"
	"os"
	"fmt"
)

func main() {
	var h handler
	h.config.fillFromYML("conf.yml")
	fmt.Println(h.config)

	http.HandleFunc("/", h.main)
	http.HandleFunc("/thumb", h.thumb)

	osFiles, _ := os.ReadDir("static")
	for _, osFile := range osFiles {
		http.HandleFunc("/" + osFile.Name(), h.static)
	}

	fmt.Println(http.ListenAndServe(h.config.Port, nil))
}