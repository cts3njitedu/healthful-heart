package services

import (
	"fmt"
	"mime/multipart"
	"github.com/cts3njitedu/healthful-heart/utils"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
	"strconv"
)

type FileService struct {
	workFile mongorepo.IWorkfileRepository
	environmentUtil utils.IEnvironmentUtility
	fileRepository mysqlrepo.IFileRepository
	rabbitService IRabbitService
}

func NewFileService(workFile mongorepo.IWorkfileRepository,
	environmentUtil utils.IEnvironmentUtility, 
	fileRepository mysqlrepo.IFileRepository, rabbitService IRabbitService) *FileService {
	return &FileService{workFile,environmentUtil, fileRepository, rabbitService}
}

func (fileService *FileService) UploadFile(file multipart.File, fileHeader * multipart.FileHeader, cred models.Credentials) error {
	fmt.Println("File Upload Endpoint Hit");
	defer file.Close()
	var userId int64
	userId, err := strconv.ParseInt(cred.UserId,10,64)
	if err != nil {
		fmt.Println(err)
	}
	newFile := &models.WorkoutFile{
		User_Id: userId,
		File_Name: fileHeader.Filename,
		Status: models.IMPORT_IN_PROGRESS,
		Version_Nb: 1,

	}
	// helper:=int32(15000)
	// mOptions:=&options.UploadOptions{
	// 	ChunkSizeBytes: &helper,
	// }

	fmt.Printf("Saving file %s to database", fileHeader.Filename);
	fileService.fileRepository.SaveFile(newFile);

	fileSize, err := fileService.workFile.StoreWorkoutFile(file, fileHeader, *newFile)

	if err!=nil {
		fmt.Println(err)
		return err
	}
	
	fmt.Printf("Write file to DB was successful. File size: %d \n", fileSize)

	fileService.rabbitService.PushFileMetaDataToQueue(newFile)
	return nil
}



