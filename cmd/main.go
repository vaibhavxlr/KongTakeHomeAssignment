package main

import (
	"log"
	"net/http"

	"github.com/vaibhavxlr/KongTakeHomeAssignment/internal"
	dbclient "github.com/vaibhavxlr/KongTakeHomeAssignment/internal/dbClient"
)

func main() {
	dbclient.ConnectMongo()
	defer dbclient.DisconnectMongo()

	mux := http.NewServeMux()
	mux.HandleFunc("/konnectAPI/v1/listServices", internal.ListServices)
	mux.HandleFunc("/konnectAPI/v1/service/{id}", internal.ServiceDetails)

	log.Println("Konnect http server listening on :8080 ...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalln("Something went wrong while starting server, err: ", err)
	}
}
