package main

import (
	"log"
	"net/http"

	"github.com/jorgemhc-dev/academy-go-q12021/packages/server"
)

func main() {
	s := server.New()

	log.Fatal(http.ListenAndServe(":8080", s.Router()))
}