package mongorepo

import (
	"github.com/cts3njitedu/healthful-heart/models"
)

type IPageRepository interface {
	GetPage(pageType string) models.Page
}