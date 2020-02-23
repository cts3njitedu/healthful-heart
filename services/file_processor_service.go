package services


import (
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
	"github.com/cts3njitedu/healthful-heart/utils"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	"bytes"
	"github.com/360EntSecGroup-Skylar/excelize"
)


type FileProcessorService struct {
	workFile mongorepo.IWorkfileRepository
	environmentUtil utils.IEnvironmentUtility
	fileRepository mysqlrepo.IFileRepository
}

func NewFileProcessorService(workFile mongorepo.IWorkfileRepository,
	environmentUtil utils.IEnvironmentUtility, fileRepository mysqlrepo.IFileRepository) *FileProcessorService {

		return &FileProcessorService{workFile, environmentUtil, fileRepository}
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

	rows := excelFile.GetRows("Sheet1")
	
    for _, row := range rows {
	
        for _, colCell := range row {
            fmt.Print(colCell, "\t")
        }
        fmt.Println()
    }
	
	return nil

}

// func printRowNum(excelFile *excelize.File) {
	
// 	rows := excelFile.GetRows("Sheet");

// 	for _, row := range rows {
// 		fmt.Println(row)
// 	}
// }