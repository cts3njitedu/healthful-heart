package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"github.com/cts3njitedu/healthful-heart/handlers"
	
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World</h1>")
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>About</h1>")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/about", about)
	r.HandleFunc("/login", handlers.Login)
	r.HandleFunc("/", index)
	http.Handle("/", r);
	fmt.Println("Server Starting...")
	

	http.ListenAndServe(GetPort(), nil)
}


func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8000"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	fmt.Println("The port is: ",port)

	return ":" + port

}