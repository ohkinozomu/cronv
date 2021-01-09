package server

import (
	"fmt"
	"net/http"
)

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	fmt.Fprintf(w, "")
}

func Serve() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/", fs)
	http.HandleFunc("/data", dataHandler)
	http.ListenAndServe(":8080", nil)
}
