package main

import (
	"net/http"
	"log"
)

func main() {
	mux := http.NewServeMux()
	
	srv := &http.Server{
		Addr:		":8080",
		Handler:	mux,
	}
	log.Fatal(srv.ListenAndServe())
}