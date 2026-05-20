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

	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("/assets", http.FileServer(http.Dir(".")))

	log.Fatal(srv.ListenAndServe())
}