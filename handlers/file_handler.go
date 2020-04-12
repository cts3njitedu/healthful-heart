package handlers

import (
	"fmt"
	"net/http"
	"github.com/cts3njitedu/healthful-heart/services"
	"github.com/cts3njitedu/healthful-heart/models"
	"strconv"
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

	creds, _ := r.Context().Value("credentials").(models.Credentials);
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile");

	location, _ := strconv.ParseInt(r.FormValue("location"),10, 64)
	metaData := models.WorkoutFile{
		Location_Id : location,
	}
	if err != nil {
		fmt.Println("Error Retrieving the file")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fHandler.fileService.UploadFile(file,handler,metaData, creds)
	
	w.Write([]byte(`{"message": "In about page"}`))

}