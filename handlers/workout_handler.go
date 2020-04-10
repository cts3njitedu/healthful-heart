package handlers

import (
	"github.com/cts3njitedu/healthful-heart/services"
	"github.com/cts3njitedu/healthful-heart/models"
	"net/http"
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type WorkoutHandler struct {
	workoutService services.IWorkoutService
}

type IWorkoutHandler interface {
	GetWorkoutDays(w http.ResponseWriter, r *http.Request)
	GetWorkoutDaysPage(w http.ResponseWriter, r *http.Request)
}

func NewWorkoutHandler(workoutService services.IWorkoutService) *WorkoutHandler {
	return &WorkoutHandler{workoutService}
}

func (handler *WorkoutHandler) GetWorkoutDays(w http.ResponseWriter, r *http.Request) {
	credentials:=context.Get(r,"credentials")
	var creds models.Credentials
	if c, ok := credentials.(models.Credentials); ok {
		creds = c
	}
	queryParam := models.QueryParams{}
	workDays, _ := handler.workoutService.GetWorkoutDays(queryParam, creds)
	json.NewEncoder(w).Encode(workDays)
}

func (handler *WorkoutHandler) GetWorkoutDaysPage(w http.ResponseWriter, r *http.Request) {
	credentials:=context.Get(r,"credentials")
	var creds models.Credentials
	if c, ok := credentials.(models.Credentials); ok {
		creds = c
	}
	
	vars := mux.Vars(r)

	queryParam := models.QueryParams{}
	queryParam.Date = vars["date"]
	heartResponse, _ := handler.workoutService.GetWorkoutDaysPage(queryParam, creds)
	json.NewEncoder(w).Encode(heartResponse)

}