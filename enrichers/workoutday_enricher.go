package enrichers

import (
	"github.com/cts3njitedu/healthful-heart/models"

)

type WorkoutDayEnricher struct {
	
}


func NewWorkoutDayEnricher() *WorkoutDayEnricher {
	return &WorkoutDayEnricher{}
}


func(enrich * WorkoutDayEnricher) Enrich(page *models.Page, pageTemplate models.Page) {

}