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
	fmt.Println("Why aren't you coming here")
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

	cell := excelFile.GetCellValue("Sheet1", "A1");

	// if err !=nil {
	// 	fmt.Println(err)
	// 	return err;
	// }
	fmt.Printf("Cell value is: %v", cell);
	return nil

}