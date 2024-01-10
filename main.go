package main

import (
	"log"
	"net/http"

	"go.chat/controllers"
)

func ControllerSetup() {
	http.HandleFunc("/", controllers.HomeHandler)
	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/message", controllers.AddMessageHandler)
	http.HandleFunc("/message/delete", controllers.DeleteMessageHandler)
}

func main() {
	log.Println("Setting up controllers")
	ControllerSetup()
	log.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
