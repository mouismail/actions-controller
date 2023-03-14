package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"

	"github.tools.sap/actions-rollout-app/utils"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/version", utils.VersionHandler)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "3000"
	}

	log.Println("Starting server on port", httpPort)

	err := http.ListenAndServe(":"+httpPort, r)

	if err != nil {
		log.Fatalf("error occured during starting the app %s", err)
	}
}
