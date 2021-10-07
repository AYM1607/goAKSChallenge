package main

import (
	"log"

	"github.com/AYM1607/goAKSChallenge/internal/server"
)

func main() {
	srvr, err := server.NewServer(":8888")
	if err != nil {
		log.Fatal("Server could not be created")
	}
	log.Fatal(srvr.ListenAndServe())
}
