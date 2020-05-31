package handlers

import (
	"github.com/cts3njitedu/healthful-heart/services"
	// "github.com/cts3njitedu/healthful-heart/models"
	"net/http"
	// "encoding/json"
	// "github.com/gorilla/mux"
	"fmt"
)

type WorkoutTypeHandler struct {
	workoutTypeService services.IWorkoutTypeService
}

type IWorkoutTypeHandler interface {
	GetCategories(w http.ResponseWriter, r *http.Request)
}

func NewWorkoutTypeHandler(workoutTypeService services.IWorkoutTypeService) *WorkoutTypeHandler {
	return &WorkoutTypeHandler{workoutTypeService}
}

func (handler *WorkoutTypeHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	c := handler.workoutTypeService.GetWorkoutTypeByIds([]int64{1,5,7});
	fmt.Printf("Categores: %+v\n", c)
}