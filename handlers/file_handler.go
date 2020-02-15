package handlers

import (
	"fmt"
	"net/http"
	"github.com/cts3njitedu/healthful-heart/services"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/gorilla/context"
)

type FileHandler struct {
	fileService services.IFileService
}

type IFileHandler interface {
	UploadFile(w http.ResponseWriter, r *http.Request)
}

func NewFileHandler(fileService services.IFileService) *FileHandler {
	return &FileHandler{fileService}
}
func (fHandler *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit");

	credentials:=context.Get(r,"credentials")
	var creds models.Credentials

	if c, ok := credentials.(models.Credentials); ok {
		creds = c
	}
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile");


	if err != nil {
		fmt.Println("Error Retrieving the file")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fHandler.fileService.UploadFile(file,handler,creds)
	
	w.Write([]byte(`{"message": "In about page"}`))

}