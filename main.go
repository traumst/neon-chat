package main

import (
	"log"
	"net/http"

	"go.chat/controllers"
)

func ControllerSetup() {
	http.HandleFunc("/", controllers.Home)
	http.HandleFunc("/login", controllers.Login)
	//http.HandleFunc("/logout", controllers.LogoutHandler)
	http.HandleFunc("/openchat", controllers.OpenChat)
	http.HandleFunc("/newchat", controllers.AddChat)
	http.HandleFunc("/message", controllers.AddMessage)
	http.HandleFunc("/message/delete", controllers.DeleteMessage)
}

func main() {
	log.Println("Setting up controllers")
	ControllerSetup()
	log.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
