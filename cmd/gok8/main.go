package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avalchev94/gok8/internal/diagnostics"
	"github.com/gorilla/mux"
)

type serverConf struct {
	port   string
	router http.Handler
	name   string
}

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

	possibleErrors := make(chan error, 2)

	configurations := []serverConf{
		{
			port:   blPort,
			router: router,
			name:   "application server",
		},
		{
			port:   diagPort,
			router: diagnostics.NewDiagnostics(),
			name:   "diagnostics server",
		},
	}

	servers := make([]*http.Server, 2)

	for i, c := range configurations {
		go func(conf serverConf, i int) {
			log.Printf("The %s is preparing to handle connections...", conf.name)
			servers[i] = &http.Server{
				Addr:    ":" + conf.port,
				Handler: conf.router,
			}
			err := servers[i].ListenAndServe()
			if err != nil {
				possibleErrors <- err
			}
		}(c, i)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-possibleErrors:
		log.Printf("Got an error %v", err)
	case sig := <-interrupt:
		log.Printf("Recieved the signal %v", sig)
	}

	for _, s := range servers {
		timeout := 5 * time.Second
		log.Printf("Shutdown with timeout: %s\n", timeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		customError := s.Shutdown(ctx)
		if customError != nil {
			fmt.Println(customError)
		}
		log.Println("Graceful shutdown")
	}
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello handler called.")
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}
