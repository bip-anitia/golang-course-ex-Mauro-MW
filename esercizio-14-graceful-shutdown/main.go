package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// TODO: Implementare HTTP server con graceful shutdown

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
