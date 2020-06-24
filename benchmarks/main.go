package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"bitbucket.org/dreamplug-backend/benchmarks/handler"
)

func main() {
	http.HandleFunc("/greet/", handler.Handler)
	http.HandleFunc("/stats/",handler.WithStats(handler.Handler))
	log.Printf("listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//go-wrk -d 5 http://localhost:8080
