package mongorepo

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"mime/multipart"
	"bytes"
)

type IPageRepository interface {
	GetPage(pageType string) models.Page
}

type IWorkfileRepository interface {
	StoreWorkoutFile(file multipart.File, fileHeader * multipart.FileHeader, newFile models.WorkoutFile) (int, error)
	RetrieveFile(file *models.WorkoutFile) (bytes.Buffer, error)
}