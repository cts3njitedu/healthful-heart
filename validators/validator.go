package validators

import (
	"github.com/cts3njitedu/healthful-heart/models"
)

type IFieldValidator interface {
	MandatoryFieldValidator(field *models.Field, v *models.Validation)
	RegexValueValidator(field *models.Field, v *models.Validation)
	LengthValidator(field *models.Field, v *models.Validation)
}

type IPageValidator interface {
	ValidatePage(page *models.Page) (error)
}