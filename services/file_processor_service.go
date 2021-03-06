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
	pageRepository mongorepo.IPageRepository
}

func NewFileProcessorService(workFile mongorepo.IWorkfileRepository,
	fileRepository mysqlrepo.IFileRepository,
	workoutTypeService IWorkoutTypeService,
	groupParserService IGroupParserService, workoutDayRepository mysqlrepo.IWorkoutDayRepository, pageRepository mongorepo.IPageRepository) *FileProcessorService {

		return &FileProcessorService{workFile, fileRepository, workoutTypeService, groupParserService, workoutDayRepository, pageRepository}
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
	workoutDays := make([]models.WorkoutDay, 0, len(workoutDayMap))
	for w := range workoutDayMap {
		workoutDay := workoutDayMap[w]
		workoutDay.Workout_File_Id = &newFile.Workout_File_Id;
		workoutDays = append(workoutDays, *workoutDay)
	}

	process.workoutDayRepository.UpdateAllWorkoutDay(workoutDays, nil)
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
	var workoutType int64
	var workoutName string
	categoryNameToCd, _ := process.workoutTypeService.GetCategories();
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
				if categoryCd, ok := categoryNameToCd[cellValue]; ok {
					categoryCdState = categoryCd
					fmt.Printf("Category Code: %v\n", categoryCdState)
				} else {
					wkType, err := process.workoutTypeService.GetWorkoutType(categoryCdState, cellValue) 
					workoutType = wkType.Workout_Type_Id;
					if workoutType == 0 || err != nil {
						fmt.Printf("Missing work type code: %v\n", cellValue)
					}
					workoutName = cellValue
					fmt.Printf("Category: %v WorkoutType: %v\n", categoryCdState, workoutType)
				}
			} else if r>0 && c>0 {
				if len(cellValue) > 0 {
					workoutDay := workoutMap[c]
					if workoutDay == nil {
						continue;
					}
					groups := process.groupParserService.GetGroups(workoutName, cellValue, categoryCdState)
					newWorkout := models.Workout {
						Workout_Type_Id: workoutType,
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