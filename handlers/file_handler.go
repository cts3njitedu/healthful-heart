package handlers

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"github.com/cts3njitedu/healthful-heart/connections"
)

type FileHandler struct {
	connection connections.IMongoConnection
}

type IFileHandler interface {
	UploadFile(w http.ResponseWriter, r *http.Request)
}

func NewFileHandler(connection connections.IMongoConnection) *FileHandler {
	return &FileHandler{connection}
}
func (fHandler *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit");

	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("myFile");


	if err != nil {
		fmt.Println("Error Retrieving the file")
		fmt.Println(err)
		return
	}
	defer file.Close()

	client,err:=fHandler.connection.GetConnection();

	if err != nil {
		fmt.Println(err)
	}

	data, err := ioutil.ReadAll(file)
	
	if err != nil {
		fmt.Println(err)
	}

	db := client.Database("WorkoutFiles")

	bucket, err := gridfs.NewBucket(db)

	if err != nil {
		fmt.Println(err)
	}

	
	uploadStream, err := bucket.OpenUploadStream("myNewFile")

	defer uploadStream.Close()

	if err != nil {
		fmt.Println(err)
	}

	fileSize, err := uploadStream.Write(data)

	if err!=nil {
		fmt.Println(err)
	}
	w.Write([]byte(`{"message": "In about page"}`))
	
	fmt.Printf("Write file to DB was successful. File size: %d \n", fileSize)
}