package main

import (
	"log"
	"net/http"

	"go.chat/src"
)

func main() {
	log.Println("Router setup...")
	src.ControllerSetup()
	log.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
