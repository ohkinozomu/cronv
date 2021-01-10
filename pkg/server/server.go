package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ohkinozomu/cronv"
)

func buildDataHandler(ctx *cronv.CronvCtx) func(http.ResponseWriter, *http.Request) {
	j, err := json.Marshal(ctx)
	if err != nil {
		log.Println(err)
	}
	return func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, string(j)) }
}

func Serve(ctx *cronv.CronvCtx) {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/", fs)
	http.HandleFunc("/data", buildDataHandler(ctx))
	http.ListenAndServe(":8080", nil)
}
