package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/ohkinozomu/cronv"
)

func buildDataHandler(ctx *cronv.CronvCtx) func(http.ResponseWriter, *http.Request) {
	j, err := json.Marshal(ctx)
	if err != nil {
		log.Println(err)
	}
	return func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, string(j)) }
}

func jsEscapeStringHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("v")
	fmt.Fprintf(w, template.JSEscapeString(strings.TrimSpace(v)))
}

func Serve(ctx *cronv.CronvCtx) {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/", fs)
	http.HandleFunc("/data", buildDataHandler(ctx))
	http.HandleFunc("/jsescapestring", jsEscapeStringHandler)
	http.ListenAndServe(":8080", nil)
}
