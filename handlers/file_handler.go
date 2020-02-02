package handlers

import (
	"fmt"
	"net/http"
	"github.com/cts3njitedu/healthful-heart/services"
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

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile");


	if err != nil {
		fmt.Println("Error Retrieving the file")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fHandler.fileService.UploadFile(file,handler)
	
	w.Write([]byte(`{"message": "In about page"}`))

}