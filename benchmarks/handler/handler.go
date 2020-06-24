package handler

import (
	"fmt"
	"log"
	"net/http"
	_"net/http/pprof"
	"bitbucket.org/dreamplug-backend/benchmarks/internal"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	msg := "hello, stranger"
	splits:=strings.Split(r.URL.Path,"/")
	email:=splits[len(splits)-1]
	if name, ok := internal.IsGopher(email); ok {
		msg = "hello gopher, " + name
	}
	_, err := fmt.Fprintln(w, msg)
	if err != nil {
		log.Printf("could not print message: %v", err)
	}
}
