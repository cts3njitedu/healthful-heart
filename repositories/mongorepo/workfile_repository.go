package mongorepo


import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/utils"
	"fmt"
	"mime/multipart"
	"io/ioutil"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"bytes"
)

type WorkfileRepository struct {
	connection connections.IMongoConnection
	environmentUtil utils.IEnvironmentUtility
}

func NewWorkfileRepository(connection connections.IMongoConnection,environmentUtil utils.IEnvironmentUtility) *WorkfileRepository {
		return &WorkfileRepository{connection, environmentUtil}
}

func (workRepo WorkfileRepository) StoreWorkoutFile(file multipart.File, fileHeader * multipart.FileHeader, newFile models.WorkoutFile) (int, error) {
	fmt.Printf("Saving file %s to mongo", fileHeader.Filename);
	
	client,err:=workRepo.connection.GetFileConnection();

	if err != nil {
		fmt.Println(err)
		return 0,err;
	}

	data, err := ioutil.ReadAll(file)
	
	if err != nil {
		fmt.Println(err)
		return 0,err;
	}
	dbName:=workRepo.environmentUtil.GetEnvironmentString("MONGODB_HEALTH_FILE_DB")

	db := client.Database(dbName)

	bucket, err := gridfs.NewBucket(db)

	if err != nil {
		fmt.Println(err)
		return 0,err;
	}

	uploadStream, err := bucket.OpenUploadStreamWithID(newFile.Workout_File_Id, fileHeader.Filename,)

	defer uploadStream.Close()

	if err != nil {
		fmt.Println(err)
		return 0,err
	}

	fileSize, err := uploadStream.Write(data)
	
	return fileSize, err


}


func (workRepo WorkfileRepository) RetrieveFile(file *models.WorkoutFile) (bytes.Buffer, error) {
	client,err:=workRepo.connection.GetFileConnection();

	if err != nil {
		fmt.Println(err)
		return bytes.Buffer{},err;
	}

	dbName:=workRepo.environmentUtil.GetEnvironmentString("MONGODB_HEALTH_FILE_DB")
	db := client.Database(dbName)

	bucket, _ := gridfs.NewBucket(db,)

	var buf bytes.Buffer

	dStream, err := bucket.DownloadToStream(file.Workout_File_Id, &buf) 

	if err != nil {
		fmt.Println(err)
		return bytes.Buffer{}, err
	}

	fmt.Printf("File size to download: %v \n", dStream)
	return buf, nil

}
