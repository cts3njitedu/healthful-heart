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
	WorkoutDaysActions(w http.ResponseWriter, r *http.Request)
	WorkoutActions(w http.ResponseWriter, r *http.Request)
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

func (handler *WorkoutHandler) WorkoutDaysActions(w http.ResponseWriter, r *http.Request) {
	
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
	fmt.Printf("Heart Request: %+v\n", heartRequest.ActionType)
	
	if heartRequest.ActionType == "VIEW_WORKOUTDATE_LOCATIONS" || heartRequest.ActionType == "VIEW_NON_WORKOUTDATE_LOCATIONS" {
		heartResponse, _ := handler.workoutService.GetWorkoutDaysLocationsView(heartRequest, creds)
		json.NewEncoder(w).Encode(heartResponse)
	} else if heartRequest.ActionType == "ADD_WORKOUTDATE_LOCATION" {
		heartResponse, _ := handler.workoutService.AddWorkoutDateLocation(heartRequest, creds)
		json.NewEncoder(w).Encode(heartResponse)
	}
	

	

}

func (handler *WorkoutHandler) WorkoutActions(w http.ResponseWriter, r *http.Request) {
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
	heartRequest.LocationId = vars["locationId"]
	fmt.Printf("Heart Request: %+v\n", heartRequest.ActionType)

	if heartRequest.ActionType == "VIEW_WORKOUTS_HEADER" {
		heartResponse, _ := handler.workoutService.GetWorkoutPageHeader(heartRequest, creds)
		json.NewEncoder(w).Encode(heartResponse)
	} else if heartRequest.ActionType == "VIEW_WORKOUTS" {
		heartResponse, _ := handler.workoutService.GetWorkouts(heartRequest, creds)
		json.NewEncoder(w).Encode(heartResponse)
	} else if heartRequest.ActionType == "VIEW_WORKOUT_DETAILS_META_INFO" {
		heartResponse, _ := handler.workoutService.GetWorkoutDetailsMetaInfo(heartRequest, creds)
		json.NewEncoder(w).Encode(heartResponse)
	} else if heartRequest.ActionType == "VIEW_WORKOUT_DETAILS" {
		heartResponse, _ := handler.workoutService.GetWorkoutDetails(heartRequest, creds)
		json.NewEncoder(w).Encode(heartResponse)
	}
}