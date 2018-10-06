package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/avalchev94/gok8/internal/diagnostics"
	"github.com/gorilla/mux"
)

func main() {
	log.Printf("Starting the application...")

	blPort := os.Getenv("PORT")
	if len(blPort) == 0 {
		log.Fatal("The port should be set")
	}
	diagPort := os.Getenv("DIAG_PORT")
	if len(diagPort) == 0 {
		log.Fatal("The diagnostics port should be set")
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handleHello)

	go func() {
		err := http.ListenAndServe(":"+blPort, router)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := http.ListenAndServe(":"+diagPort, diagnostics.NewDiagnostics())
	if err != nil {
		log.Fatal(err)
	}

}

func handleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}
