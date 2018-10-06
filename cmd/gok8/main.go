package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/avalchev94/gok8/internal/diagnostics"
	"github.com/gorilla/mux"
)

func main() {
	log.Printf("Hello, World")

	router := mux.NewRouter()
	router.HandleFunc("/", handleHello)

	go func() {
		err := http.ListenAndServe(":8080", router)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := http.ListenAndServe(":8585", diagnostics.NewDiagnostics())
	if err != nil {
		log.Fatal(err)
	}

}

func handleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}
