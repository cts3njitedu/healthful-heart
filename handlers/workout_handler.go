package handlers

import (
	"github.com/cts3njitedu/healthful-heart/services"
	"github.com/cts3njitedu/healthful-heart/models"
	"net/http"
	"encoding/json"
	"github.com/gorilla/context"
)

type WorkoutHandler struct {
	workoutService services.IWorkoutService
}

type IWorkoutHandler interface {
	GetWorkoutDays(w http.ResponseWriter, r *http.Request)
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