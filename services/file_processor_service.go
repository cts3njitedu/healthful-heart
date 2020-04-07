package services


import (
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	"bytes"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strings"
	"time"
)


type FileProcessorService struct {
	workFile mongorepo.IWorkfileRepository
	fileRepository mysqlrepo.IFileRepository
	workoutTypeService IWorkoutTypeService
	groupParserService IGroupParserService
	workoutDayRepository mysqlrepo.IWorkoutDayRepository
}

func NewFileProcessorService(workFile mongorepo.IWorkfileRepository,
	fileRepository mysqlrepo.IFileRepository,
	workoutTypeService IWorkoutTypeService,
	groupParserService IGroupParserService, workoutDayRepository mysqlrepo.IWorkoutDayRepository) *FileProcessorService {

		return &FileProcessorService{workFile, fileRepository, workoutTypeService, groupParserService, workoutDayRepository}
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

	workoutDayMap := ConvertFileToWorkoutDay(excelFile, newFile, process)
	
	for w := range workoutDayMap {
		workoutDay := workoutDayMap[w]
		workoutDay.Workout_File_Id = newFile.Workout_File_Id;
		process.workoutDayRepository.SaveWorkoutDay(workoutDay)
	}

	_, err = process.fileRepository.UpdateFileStatus(&newFile, models.FILE_PROCESSED)

	if err != nil {
		fmt.Println(err)
		return err;
	}

	return nil

}

func ConvertFileToWorkoutDay(excelFile *excelize.File, metaData models.WorkoutFile, process *FileProcessorService) (map[int]*models.WorkoutDay) {
	rows := excelFile.GetRows("Sheet1")
	workoutMap := make(map[int]*models.WorkoutDay)
	var categoryCdState string
	var workoutType string
	var workoutName string
    for r, row := range rows {
	
        for c, colCell := range row {
			cellValue := strings.TrimSpace(colCell)
			if len(cellValue) == 0 {
				continue;
			}
            if r==0 && c>0 {
				date, err := time.Parse("1/2/2006", cellValue)
				if err != nil {
					fmt.Println(err)
					continue;
				}
				dateFormat := date.Format("2006-01-02 15:04:05")
				fmt.Printf("Date Value: %+v\n", dateFormat)
				workoutDay := &models.WorkoutDay {
					Workout_Date: dateFormat,
					User_Id: metaData.User_Id,
					Location_Id: metaData.Location_Id,
				}
				workoutMap[c] = workoutDay
			} else if r>0 && c==0 {
				fmt.Printf("Row Value: %+v, Cell Value: %+v\n",r, cellValue)
				categoryCd, err := process.workoutTypeService.GetCategoryCodeFromName(cellValue)
				if err == nil {
					categoryCdState = categoryCd
					fmt.Printf("Category Code State: %v\n", categoryCdState);
				} else {
					workoutType = process.workoutTypeService.GetWorkoutTypeCode(categoryCdState, cellValue)
					if len(workoutType) == 0 {
						fmt.Printf("Missing work type code: %v", cellValue)
					}
					workoutType = strings.TrimSpace(workoutType)
					workoutName = cellValue
					fmt.Printf("Category: %v WorkoutType: %v\n", categoryCdState, workoutType)
				}
			} else if r>0 && c>0 {
				if len(cellValue) > 0 {
					workoutDay := workoutMap[c]
					if workoutDay == nil {
						continue;
					}
					groups := process.groupParserService.GetGroups(workoutType, cellValue, categoryCdState)
					newWorkout := models.Workout {
						Workout_Type_Cd: workoutType,
						Workout_Name: workoutName,
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