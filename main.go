package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"github.com/cts3njitedu/healthful-heart/handlers"
	"github.com/justinas/alice"
	"github.com/cts3njitedu/healthful-heart/factories"
	
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World</h1>")
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>About</h1>")
}

func main() {
	r := mux.NewRouter()
	getLoginPage := http.HandlerFunc(factories.GetLoginHandler().GetLoginPage);
	postLoginPage := factories.GetLoginHandler().PostLoginPage;
	getSignupPage := http.HandlerFunc(factories.GetSignupHandler().GetSignUpPage);
	postSignupPage := factories.GetSignupHandler().PostSignUpPage;
	getAboutPage := http.HandlerFunc(factories.GetAboutHandler().GetAboutPage);
	tokenHandler := http.HandlerFunc(factories.GetTokenHandler().GetToken);
	validateTokenHandler := factories.GetTokenHandler().ValidateToken;
	uploadFile := http.HandlerFunc(factories.GetFileHandler().UploadFile);
	getWorkoutDays := http.HandlerFunc(factories.GetWorkoutHandler().GetWorkoutDays);
	workType := http.HandlerFunc(factories.GetWorkoutTypeHandler().GetCategories)
	preChain := alice.New(handlers.Logging, validateTokenHandler);
	workoutDaysRoutes(r, "/workoutDays", preChain)
	r.Handle("/upload", alice.New(handlers.Logging, validateTokenHandler).Then(uploadFile)).Methods("POST")
	r.Handle("/about", alice.New(handlers.Logging, validateTokenHandler).Then(getAboutPage)).Methods("GET")
	r.Handle("/calendar", alice.New(handlers.Logging, validateTokenHandler).Then(getWorkoutDays)).Methods("GET")
	r.Handle("/login", alice.New(handlers.Logging).Then(getLoginPage)).Methods("GET");
	r.Handle("/login", alice.New(handlers.Logging, postLoginPage).Then(tokenHandler)).Methods("POST");
	r.Handle("/signup", alice.New(handlers.Logging).Then(getSignupPage)).Methods("GET");
	r.Handle("/signup", alice.New(handlers.Logging,postSignupPage).Then(tokenHandler)).Methods("POST");
	r.Handle("/workoutType", alice.New(handlers.Logging).Then(workType)).Methods("GET","POST")
	r.HandleFunc("/", index)
	http.Handle("/", r);
	fmt.Println("Server Starting...")

	http.ListenAndServe(GetPort(), nil)
}


func workoutDaysRoutes(r *mux.Router, prefix string, chain alice.Chain) {
	sub := mux.NewRouter()
	validateTokenHandler := factories.GetTokenHandler().ValidateToken;
	preChain := alice.New(handlers.Logging, validateTokenHandler);
	sub.HandleFunc(prefix ,factories.GetWorkoutHandler().WorkoutDaysActions).Methods("GET", "POST");
	sub.HandleFunc(prefix + "/", factories.GetWorkoutHandler().WorkoutActions).Methods("GET", "POST")
	sub.HandleFunc(prefix + "/{date}" ,factories.GetWorkoutHandler().WorkoutDaysActions).Methods("GET", "POST");
	sub.HandleFunc(prefix + "/{date}/locations", factories.GetWorkoutHandler().WorkoutDaysActions).Methods("POST");
	sub.HandleFunc(prefix + "/{date}/locations/{locationId}", factories.GetWorkoutHandler().WorkoutActions).Methods("GET", "POST")
	r.Handle(prefix, preChain.Then(sub))
	r.Handle(prefix + "/{path:.*}", preChain.Then(sub));
}

// func workoutRoutes(r *mux.Router, prefix string, chain alice.Chain) {
// 	sub := mux.NewRouter()
// 	validateTokenHandler := factories.GetTokenHandler().ValidateToken;
// 	preChain := alice.New(handlers.Logging, validateTokenHandler);
// 	sub.HandleFunc(prefix, factories.GetWorkoutHandler().WorkoutActions).Methods("GET", "POST")
// 	sub.HandleFunc(prefix + "/", factories.GetWorkoutHandler().WorkoutActions).Methods("GET", "POST")
// 	sub.HandleFunc(prefix + "/headers", factories.GetWorkoutHandler().WorkoutActions).Methods("GET", "POST")
// 	r.Handle(prefix, preChain.Then(sub))
// 	r.Handle(prefix + "/{path:.*}", preChain.Then(sub));

// }

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