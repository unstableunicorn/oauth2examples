package main

import (
	"log"
	"net/http"

	"github.com/unstableunicorn/oauth2examples/server/oauth2"
)

func main() {
	handler := http.HandlerFunc(oauth2.AuthServer)
	log.Fatal(http.ListenAndServe(":3002", handler))
}
