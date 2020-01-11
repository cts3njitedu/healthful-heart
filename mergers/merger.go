package mergers

import (

	"github.com/cts3njitedu/healthful-heart/models"
)


type IPageMerger interface {
	MergeRequestPageToPage(requestPage *models.Page, page models.Page)
}