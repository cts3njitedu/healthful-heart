package services_test

import (
	"testing"
	"fmt"
	"github.com/cts3njitedu/healthful-heart/services"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/stretchr/testify/assert"

)

type CategoryRepositoryMock struct {}

func (mock CategoryRepositoryMock) GetCategories() ([]models.Category, error) {
	categories := []models.Category {
		models.Category {
			Category_Cd: "AB",
			Category_Name: "Abdominal",
		},
		models.Category {
			Category_Cd: "SH",
			Category_Name: "Shoulders",
		},
	}

	return categories, nil
}
type WorkoutTypeRepositoryMock struct {}

func (mock WorkoutTypeRepositoryMock) GetWorkoutTypes() ([]models.WorkoutType, error) {
	workOutTypes := []models.WorkoutType {
		models.WorkoutType {
			Name: "Workout Type Name1",
			Category_Cd: "BC",
			Workout_Type_Cd: "BC1",
		},
		models.WorkoutType {
			Name: "Workout Type Name2",
			Category_Cd: "BK",
			Workout_Type_Cd: "BK1",
		},
		models.WorkoutType {
			Name: "Workout Type Name3",
			Category_Cd: "SH",
			Workout_Type_Cd: "SH2",
		},
		models.WorkoutType {
			Name: "Workout Type Name4",
			Category_Cd: "BC",
			Workout_Type_Cd: "BC2",
		},
		models.WorkoutType {
			Name: "Workout Type Name5",
			Category_Cd: "CH",
			Workout_Type_Cd: "CH1",
		},
		models.WorkoutType {
			Name: "Workout Type Name6",
			Category_Cd: "CH",
			Workout_Type_Cd: "CH10",
		},
		models.WorkoutType {
			Name: "Workout Type Name7",
			Category_Cd: "TR",
			Workout_Type_Cd: "TR6",
		},
	}

	return workOutTypes, nil
}


func TestGetCategoryCode(t *testing.T) {
	t.Run("Testing Get Category Code", func(t *testing.T) {
		assert := assert.New(t)
		repo := WorkoutTypeRepositoryMock{}
		catRepo := CategoryRepositoryMock{}
		service :=  services.NewWorkoutTypeService(repo, catRepo)
		fmt.Printf("Service type: %T", service)
		workoutTypeCd := service.GetWorkoutTypeCode("CH", "Workout Type Name5")
		assert.Equal("CH1", workoutTypeCd,"The workout type codes are not equal")
	})
	
}

func TestGetCategoryCodeBiceps(t *testing.T) {
	t.Run("Testing Get Category Code Biceps", func(t *testing.T) {
		assert := assert.New(t)
		repo := WorkoutTypeRepositoryMock{}
		catRepo := CategoryRepositoryMock{}
		service :=  services.NewWorkoutTypeService(repo, catRepo)
		fmt.Printf("Service type: %T", service)
		workoutTypeCd := service.GetWorkoutTypeCode("BC", "Workout Type Name1")
		assert.Equal("BC1", workoutTypeCd,"The workout type codes are not equal")
	})
	
}

func TestGetCategoryCodeWorkoutTypeCodeMissing(t *testing.T) {
	t.Run("Testing Get Category Code Workout Type Code Missing", func(t *testing.T) {
		assert := assert.New(t)
		repo := WorkoutTypeRepositoryMock{}
		catRepo := CategoryRepositoryMock{}
		service :=  services.NewWorkoutTypeService(repo, catRepo)
		fmt.Printf("Service type: %T", service)
		workoutTypeCd := service.GetWorkoutTypeCode("BC", "Workout Type Name20")
		assert.Equal("Worko", workoutTypeCd,"The workout type codes are not equal")
	})
	
}


func TestGetCategoryCodeCategoryCodeMissing(t *testing.T) {
	t.Run("Testing Get Category Code Category Code Missing", func(t *testing.T) {
		assert := assert.New(t)
		repo := WorkoutTypeRepositoryMock{}
		catRepo := CategoryRepositoryMock{}
		service :=  services.NewWorkoutTypeService(repo, catRepo)
		fmt.Printf("Service type: %T", service)
		workoutTypeCd := service.GetWorkoutTypeCode("BC8e", "Workout Type Name1")
		assert.Equal("Worko", workoutTypeCd,"The workout type codes are not equal")
	})
	
}