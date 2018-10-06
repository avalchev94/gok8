package main

import (
	"log"
	"net/http"
)

func main() {
	log.Printf("Hello, World")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
