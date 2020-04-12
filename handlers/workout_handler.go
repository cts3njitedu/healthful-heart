package handlers

import (
	"github.com/cts3njitedu/healthful-heart/services"
	"github.com/cts3njitedu/healthful-heart/models"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
)

type WorkoutHandler struct {
	workoutService services.IWorkoutService
}

type IWorkoutHandler interface {
	GetWorkoutDays(w http.ResponseWriter, r *http.Request)
	GetWorkoutDaysPage(w http.ResponseWriter, r *http.Request)
	GetWorkoutDaysLocationPage(w http.ResponseWriter, r *http.Request)
}

func NewWorkoutHandler(workoutService services.IWorkoutService) *WorkoutHandler {
	return &WorkoutHandler{workoutService}
}

func (handler *WorkoutHandler) GetWorkoutDays(w http.ResponseWriter, r *http.Request) {
	// credentials:=context.Get(r,"credentials")
	// var creds models.Credentials
	// if c, ok := credentials.(models.Credentials); ok {
	// 	creds = c
	// }
	creds, _ := r.Context().Value("credentials").(models.Credentials);
	fmt.Printf("Handler Credentials: %+v\n", creds)
	queryParam := models.QueryParams{}
	workDays, _ := handler.workoutService.GetWorkoutDays(queryParam, creds)
	json.NewEncoder(w).Encode(workDays)
}

func (handler *WorkoutHandler) GetWorkoutDaysPage(w http.ResponseWriter, r *http.Request) {
	// credentials:=context.Get(r,"credentials")
	// var creds models.Credentials
	// if c, ok := credentials.(models.Credentials); ok {
	// 	creds = c
	// }
	creds, _ := r.Context().Value("credentials").(models.Credentials);
	fmt.Printf("Handler Credentials: %+v\n", creds)
	vars := mux.Vars(r)

	queryParam := models.QueryParams{}
	queryParam.Date = vars["date"]
	heartResponse, _ := handler.workoutService.GetWorkoutDaysPage(queryParam, creds)
	json.NewEncoder(w).Encode(heartResponse)

}

func (handler *WorkoutHandler) GetWorkoutDaysLocationPage(w http.ResponseWriter, r *http.Request) {
	// credentials:=context.Get(r,"credentials")
	// var creds models.Credentials
	// if c, ok := credentials.(models.Credentials); ok {
	// 	creds = c
	// }
	creds, _ := r.Context().Value("credentials").(models.Credentials);
	fmt.Printf("Handler credential: %+v\n", creds)
	
	vars := mux.Vars(r)
	heartRequest := models.HeartRequest{};
	err := json.NewDecoder(r.Body).Decode(&heartRequest)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	heartRequest.Date = vars["date"]
	fmt.Printf("Heart Request: %+v\n", heartRequest)
	if models.VIEW_LOCATIONS == heartRequest.ActionType {
		heartResponse, _ := handler.workoutService.GetWorkoutDaysLocationsView(heartRequest, creds)
		json.NewEncoder(w).Encode(heartResponse)
	}
	

}