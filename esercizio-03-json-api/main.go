package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// TODO: Implementare il JSON API Server

	fmt.Println("Starting JSON API Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
