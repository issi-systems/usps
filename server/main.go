package main

import (
	"io"
	"log"
	"net/http"
	"usps"
)

func main() {
	usps.Init()
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, world"))
	})
	http.HandleFunc("/address", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		b, e := io.ReadAll(r.Body)
		if e != nil || len(b) == 0 {
			log.Printf("query: %v, %s", e, b)
			return
		}
		tj, e := usps.ValidateJson(b)
		if e != nil {
			log.Printf("9 %v %s", e, string(b))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(tj)
		log.Printf("response %s", tj)
	})

	log.Printf("listening on 8073")
	log.Fatal(http.ListenAndServe(":8073", nil))
}
