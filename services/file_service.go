package services

import (
	"fmt"
	"io/ioutil"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"github.com/cts3njitedu/healthful-heart/connections"
	"mime/multipart"
	"github.com/cts3njitedu/healthful-heart/utils"
)

type FileService struct {
	connection connections.IMongoConnection
	environmentUtil utils.IEnvironmentUtility
}

func NewFileService(connection connections.IMongoConnection,environmentUtil utils.IEnvironmentUtility) *FileService {
	return &FileService{connection,environmentUtil}
}

func (fileService *FileService) UploadFile(file multipart.File, fileHeader * multipart.FileHeader) error {
	fmt.Println("File Upload Endpoint Hit");

	client,err:=fileService.connection.GetConnection();

	if err != nil {
		fmt.Println(err)
	}

	data, err := ioutil.ReadAll(file)
	
	if err != nil {
		fmt.Println(err)
		return err;
	}
	dbName:=fileService.environmentUtil.GetEnvironmentString("MONGODB_HEALTH_CONFIG_DB")

	db := client.Database(dbName)

	bucket, err := gridfs.NewBucket(db)

	if err != nil {
		fmt.Println(err)
		return err;
	}

	
	uploadStream, err := bucket.OpenUploadStream(fileHeader.Filename)

	defer uploadStream.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	fileSize, err := uploadStream.Write(data)

	if err!=nil {
		fmt.Println(err)
		return err
	}
	
	fmt.Printf("Write file to DB was successful. File size: %d \n", fileSize)
	return nil
}



