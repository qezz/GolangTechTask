package main

import (
	"net/http"
	"log"
	"flag"

	"github.com/buffup/GolangTechTask/cmd/server/internal/handlers"

)

func main() {

    var port string
    flag.StringVar(&port, "p", "8080", "Specify listening port. Default is 8080")
    flag.Parse()

	log.Print("Service is started, localhost:" + port)
	// Replace with DB impl
	store := handlers.NewInMemStore()

	routes := handlers.Routes(store)
	if err := http.ListenAndServe(":" + port, routes); err != nil {
		panic(err)
	}
}
