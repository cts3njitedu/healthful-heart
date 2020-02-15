package services

import (
	"fmt"
	"io/ioutil"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"github.com/cts3njitedu/healthful-heart/connections"
	"mime/multipart"
	"github.com/cts3njitedu/healthful-heart/utils"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"strconv"
)

type FileService struct {
	connection connections.IMongoConnection
	environmentUtil utils.IEnvironmentUtility
	fileRepository mysqlrepo.IFileRepository
	rabbitService IRabbitService
}

func NewFileService(connection connections.IMongoConnection,
	environmentUtil utils.IEnvironmentUtility, 
	fileRepository mysqlrepo.IFileRepository, rabbitService IRabbitService) *FileService {
	return &FileService{connection,environmentUtil, fileRepository, rabbitService}
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

	client,err:=fileService.connection.GetFileConnection();

	if err != nil {
		fmt.Println(err)
		return err;
	}

	data, err := ioutil.ReadAll(file)
	
	if err != nil {
		fmt.Println(err)
		return err;
	}
	dbName:=fileService.environmentUtil.GetEnvironmentString("MONGODB_HEALTH_FILE_DB")

	db := client.Database(dbName)

	bucket, err := gridfs.NewBucket(db)

	if err != nil {
		fmt.Println(err)
		return err;
	}
	// helper:=int32(15000)
	// mOptions:=&options.UploadOptions{
	// 	ChunkSizeBytes: &helper,
	// }

	fmt.Printf("Saving file %s to database", fileHeader.Filename);
	fileService.fileRepository.SaveFile(newFile);
	
	uploadStream, err := bucket.OpenUploadStreamWithID(newFile.Workout_File_Id, fileHeader.Filename,)

	defer uploadStream.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Saving file %s to mongo", fileHeader.Filename);
	fileSize, err := uploadStream.Write(data)

	if err!=nil {
		fmt.Println(err)
		return err
	}

	
	fmt.Printf("Write file to DB was successful. File size: %d \n", fileSize)

	fileService.rabbitService.PushFileMetaDataToQueue(newFile)
	return nil
}



