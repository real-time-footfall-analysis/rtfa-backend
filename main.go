package main

import (
	"log"
	"net/http"
)

type TestMessage struct {
	Message string `json:"message"`
}

func main() {
	a := App{}
	initialize(&a)

	log.Fatal(http.ListenAndServe(":80", a.Router))
}
