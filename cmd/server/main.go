package main

import (
	"log"

	"github.com/AYM1607/goAKSChallenge/internal/server"
)

func main() {
	srvr := server.New(":8888")
	log.Fatal(srvr.ListenAndServe())
}
