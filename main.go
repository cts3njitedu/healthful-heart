package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"github.com/cts3njitedu/healthful-heart/handlers"
	"github.com/justinas/alice"
	
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World</h1>")
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>About</h1>")
}

func main() {
	r := mux.NewRouter()
	getLoginPage := http.HandlerFunc(handlers.GetLoginPage);
	postLoginPage := http.HandlerFunc(handlers.PostLoginPage);
	r.HandleFunc("/about", about)
	r.Handle("/login", alice.New(handlers.Logging).Then(getLoginPage)).Methods("GET");
	r.Handle("/login", alice.New(handlers.Logging).Then(postLoginPage)).Methods("POST");
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