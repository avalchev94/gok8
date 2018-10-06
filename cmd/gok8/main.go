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
		log.Print("The application server is starting...")
		err := http.ListenAndServe(":"+blPort, router)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Print("The diagnostics server is starting...")
	err := http.ListenAndServe(":"+diagPort, diagnostics.NewDiagnostics())
	if err != nil {
		log.Fatal(err)
	}

}

func handleHello(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello handler called.")
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}
