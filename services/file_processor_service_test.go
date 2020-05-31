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
	"github.com/jinzhu/gorm"
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

type PageRepositoryMock struct {}

func (repo PageRepositoryMock) GetPage(pageType string) models.Page {
	return models.Page{}
}
type WorkoutTypeServiceMock struct {}

func (serv WorkoutTypeServiceMock) GetWorkoutType(categoryCd string, workoutTypeName string) (models.WorkoutType, error) {
	return models.WorkoutType{
		Workout_Type_Id: 1,
		Workout_Type_Desc: workoutTypeName,
		Category_Cd : categoryCd,
	}, nil
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
func (serv WorkoutTypeServiceMock) GetCategories() (map[string]string, map[string]string) {
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
	n := map[string]string {
		"Abs" : "AB",
		"Biceps" : "BC",
		"Back" : "BK",
		"Chest" : "CH",
		"Default" : "DF",
		"Legs" : "LG",
		"Shoulders" : "SH",
		"Triceps" : "TR",
	}
	return m, n
}

func (serv WorkoutTypeServiceMock) GetCategoriesAndWorkoutsMap(catCode string) (map[string]map[int64]models.WorkoutType) {
	return nil
}
func (serv WorkoutTypeServiceMock) GetSortedCategoriesAndWorkoutTypes() ([]models.SortedCategoryWorkoutType) {
	return nil
}
func (serv WorkoutTypeServiceMock) GetWorkoutTypeByIds(ids []int64) (map[int64]models.WorkoutType) {
	return nil
}
type WorkoutRepositoryMock struct {}

func (repo WorkoutRepositoryMock) SaveWorkoutDay(workDay *models.WorkoutDay, tx *gorm.DB) error {
	return nil
}

func (repo WorkoutRepositoryMock) GetWorkoutDays(userId string) ([]models.WorkoutDay, error) {
	return nil, nil
}

func (repo WorkoutRepositoryMock) GetWorkoutDaysByParams(queryOptions models.QueryOptions) ([]models.WorkoutDay, error) {
	return nil, nil
}

func (repo WorkoutRepositoryMock) GetWorkoutDaysSpecifyColumns(queryOptions models.QueryOptions) ([]models.WorkoutDay, error) {
	return nil, nil
}

func(repo WorkoutRepositoryMock) SaveWorkoutDayLocation(workDay *models.WorkoutDay) (*models.WorkoutDay, error) {
	return nil, nil
}

func(repo WorkoutRepositoryMock) DeleteWorkoutDays(ids map[string][]string, tx *gorm.DB) bool {
	return true
}

func(repo WorkoutRepositoryMock) UpdateAllWorkoutDay(workDays []models.WorkoutDay, ids map[string][]string) error {
	return nil
}

func (repo WorkoutRepositoryMock) GetWorkoutDaysLocationByParams(queryOptions models.QueryOptions) ([]models.WorkoutDay, error) {
	return nil, nil
}
func TestConvertFileToWorkDayMap(t *testing.T) {
	t.Run("Testing Convert File to WorkDay Map", func(t *testing.T) {
		assert := assert.New(t)
		workFileRepo := WorkfileRepositoryMock{};
		workTypeServ := WorkoutTypeServiceMock{};
		fileMock := FileRepositoryMock{}
		workDayRepo := WorkoutRepositoryMock{}
		service :=  services.NewGroupParserService()
		pageRepo := PageRepositoryMock{}
		fileProcServ := services.NewFileProcessorService(workFileRepo,fileMock,workTypeServ, service, workDayRepo, pageRepo)
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