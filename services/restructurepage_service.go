package services

import (

	"github.com/cts3njitedu/healthful-heart/models"
)

type RestructurePageService struct {}

func NewRestructurePageService() *RestructurePageService {
	return &RestructurePageService{}
}
func (restruct *RestructurePageService) RestructureLoginPage(page *models.Page) {
	for s := range page.Sections {
		var section = &page.Sections[s];
		for f := range section.Fields {
			var field  = &section.Fields[f]
			field.IsHidden=false
		}
	}
}