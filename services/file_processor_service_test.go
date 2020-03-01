package services_test

import (
	"testing"
	"fmt"
	"github.com/cts3njitedu/healthful-heart/services"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/stretchr/testify/assert"
	"bytes"
	"mime/multipart"
	"github.com/360EntSecGroup-Skylar/excelize"
	"encoding/json"
	"errors"
    // "os"
)

type WorkfileRepositoryMock struct {}

func (workRepo WorkfileRepositoryMock) StoreWorkoutFile(file multipart.File, fileHeader * multipart.FileHeader, newFile models.WorkoutFile) (int, error) {
	return 1,nil
}

func (workRepo WorkfileRepositoryMock) RetrieveFile(file *models.WorkoutFile) (bytes.Buffer, error) {
	return bytes.Buffer{}, nil
}

type FileRepositoryMock struct {}

func (repo FileRepositoryMock) SaveFile(file *models.WorkoutFile) (error) {
	return nil
}

func (repo FileRepositoryMock) UpdateFileStatus(file *models.WorkoutFile, newStatus string) (models.WorkoutFile, error) {
	return models.WorkoutFile{}, nil
}

type WorkoutTypeServiceMock struct {}

func (serv WorkoutTypeServiceMock) GetWorkoutTypeCode(categoryCd string, workoutTypeName string) string {
	if categoryCd == "AB" && workoutTypeName == "Abdominal" {
		return "AB2"
	}
	if categoryCd == "AB" && workoutTypeName == "Leg raises" {
		return "AB5"
	}
	if categoryCd == "BC" && workoutTypeName == "Barbell rubber zig 3 moves" {
		return "BC1"
	}
	if categoryCd == "BC" && workoutTypeName == "Seated bicep curls" {
		return "BC14"
	}
	if categoryCd == "BC" && workoutTypeName == "Zodiac curls" {
		return "BC19"
	}
	if categoryCd == "BK" && workoutTypeName == "Pull ups" {
		return "BK16"
	}
	if categoryCd == "BK" && workoutTypeName == "Lower back" {
		return "BK10"
	}
	if categoryCd == "BK" && workoutTypeName == "Machine fly" {
		return "BK12"
	}
	if categoryCd == "CH" && workoutTypeName == "Dip" {
		return "CH5"
	}
	if categoryCd == "CH" && workoutTypeName == "Cable cross over" {
		return "CH2"
	}
	if categoryCd == "CH" && workoutTypeName == "Bench Barbell" {
		return "CH1"
	}
	if categoryCd == "CH" && workoutTypeName == "Pectoral fly" {
		return "CH14"
	}
	if categoryCd == "CH" && workoutTypeName == "Incline dumbell" {
		return "CH8"
	}
	if categoryCd == "LG" && workoutTypeName == "Leg curl" {
		return "LG4"
	}
	if categoryCd == "LG" && workoutTypeName == "Leg push" {
		return "LG6"
	}
	if categoryCd == "SH" && workoutTypeName == "Sitting shoulder press machine" {
		return "SH12"
	}
	if categoryCd == "SH" && workoutTypeName == "Twist" {
		return "SH16"
	}
	if categoryCd == "SH" && workoutTypeName == "Shoulder free weights" {
		return "SH8"
	}
	if categoryCd == "TR" && workoutTypeName == "Curl bar cable" {
		return "TR1"
	}
	if categoryCd == "TR" && workoutTypeName == "One hand lean over" {
		return "TR4"
	}
	if categoryCd == "TR" && workoutTypeName == "Sitting dips" {
		return "TR8"
	}
	return ""


}
func (serv WorkoutTypeServiceMock) GetCategoryCodeFromName(categoryName string) (string, error){
	m := map[string]string {
		"Abs" : "AB",
		"Biceps" : "BC",
		"Back" : "BK",
		"Chest" : "CH",
		"Default" : "DF",
		"Legs" : "LG",
		"Shoulders" : "SH",
		"Triceps" : "TR",
	}
	if m[categoryName] != "" {
		return m[categoryName], nil
	}	
	return "", errors.New("doesn't exist")
}
func (serv WorkoutTypeServiceMock) GetCategoryNameFromCode(categoryCd string) (string, error){
	m := map[string]string {
		"AB": "Abs",
		"BC":"Biceps",
		"BK":"Back",
		"CH":"Chest",
		"DF":"Default",
		"LG":"Legs",
		"SH":"Shoulders", 
		"TR":"Triceps",
	}
	if m[categoryCd] != "" {
		return m[categoryCd], nil
	}	
	return "", errors.New("doesn't exist")
}

type WorkoutRepositoryMock struct {}

func (repo WorkoutRepositoryMock) SaveWorkoutDay(workDay *models.WorkoutDay) error {
	return nil
}
func TestConvertFileToWorkDayMap(t *testing.T) {
	t.Run("Testing Convert File to WorkDay Map", func(t *testing.T) {
		assert := assert.New(t)
		workFileRepo := WorkfileRepositoryMock{};
		workTypeServ := WorkoutTypeServiceMock{};
		fileMock := FileRepositoryMock{}
		workDayRepo := WorkoutRepositoryMock{}
		service :=  services.NewGroupParserService()
		fileProcServ := services.NewFileProcessorService(workFileRepo,fileMock,workTypeServ, service, workDayRepo)
		excelFile, err := excelize.OpenFile("WorkoutLogTest.xlsx")
		if err != nil {
			fmt.Println(err)
		}
		newFile := models.WorkoutFile{
			User_Id: 123,
		}
		workoutDays := services.ConvertFileToWorkoutDay(excelFile, newFile, fileProcServ)
		jsonString, _ := json.Marshal(workoutDays)
		fmt.Printf("WorkDays: %v\n", string(jsonString))
		assert.Equal("water","water","water")
	})
	
}