package services


import (
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	"bytes"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strings"
)


type FileProcessorService struct {
	workFile mongorepo.IWorkfileRepository
	fileRepository mysqlrepo.IFileRepository
	workoutTypeService IWorkoutTypeService
	groupParserService IGroupParserService
}

func NewFileProcessorService(workFile mongorepo.IWorkfileRepository,
	fileRepository mysqlrepo.IFileRepository,
	workoutTypeService IWorkoutTypeService,groupParserService IGroupParserService) *FileProcessorService {

		return &FileProcessorService{workFile, fileRepository, workoutTypeService, groupParserService}
}

func (process *FileProcessorService) ProcessWorkoutFile(file models.WorkoutFile) (error) {

	fmt.Println("Processing file....")
	newFile, err := process.fileRepository.UpdateFileStatus(&file, models.PROCESSING_FILE)

	if err!=nil {
		fmt.Println(err)
		return err;
	}
	
	data, err := process.workFile.RetrieveFile(&newFile)

	if err != nil {
		fmt.Println(err)
		return err;
	}

	fileReader := bytes.NewReader(data.Bytes())

	excelFile, err := excelize.OpenReader(fileReader)

	if err != nil {
		fmt.Println(err)
		return err;
	}

	ConvertFileToWorkoutDay(excelFile, process)
	
	return nil

}

func ConvertFileToWorkoutDay(excelFile *excelize.File, process *FileProcessorService) (map[int]*models.WorkoutDay) {
	rows := excelFile.GetRows("Sheet1")
	workoutMap := make(map[int]*models.WorkoutDay)
	var categoryCdState string
	var workoutType string
    for r, row := range rows {
	
        for c, colCell := range row {
			cellValue := strings.TrimSpace(colCell)
            if r==0 && c>0 {
				workoutDay := &models.WorkoutDay {
					Workout_Date: cellValue,
				}
				workoutMap[c] = workoutDay
			} else if r>0 && c==0 {
				categoryCd, err := process.workoutTypeService.GetCategoryCodeFromName(cellValue)
				if err == nil {
					categoryCdState = categoryCd
				} else {
					workoutType = process.workoutTypeService.GetWorkoutTypeCode(categoryCdState, cellValue)
					workoutType = strings.TrimSpace(workoutType)
				}
			} else if r>0 && c>0 {
				if len(cellValue) > 0 {
					workoutDay := workoutMap[c]
					groups := process.groupParserService.GetGroups(workoutType, cellValue, categoryCdState)
					newWorkout := models.Workout {
						Workout_Type_Cd: workoutType,
						Groups: groups,
					}
					if workoutDay.Workouts != nil {
						workoutDay.Workouts = append(workoutDay.Workouts, newWorkout)
					} else {
						newWorkouts := make([]models.Workout, 0, 10)
						newWorkouts = append(newWorkouts, newWorkout)
						workoutDay.Workouts = newWorkouts
					}
				}
			}
        }
	}
	
	return workoutMap
}