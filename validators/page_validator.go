package validators

import (

	"github.com/cts3njitedu/healthful-heart/models"
)

type PageValidator struct {
	fieldValidator IFieldValidator
}

type ValidationError struct {
	s string
}
func NewPageValidator(fieldValidator IFieldValidator) *PageValidator {
	return &PageValidator{fieldValidator}
}
func (v * ValidationError) Error() string {
	return v.s
}

func (pageValidator PageValidator) ValidatePage(page *models.Page) (error) {
	var isPageInvalid bool = false
	for s := range page.Sections {
		section := &page.Sections[s];
		for f:= range section.Fields {
			field := &section.Fields[f]
			for v := range field.Validations {
				validation := &field.Validations[v]
				switch validation.ValidationName {
				case "MANDATORY":
					pageValidator.fieldValidator.MandatoryFieldValidator(field, validation)
				case "REGEX":
					pageValidator.fieldValidator.RegexValueValidator(field, validation)
				case "LENGTH":
					pageValidator.fieldValidator.LengthValidator(field, validation)

				}
			}
			isPageInvalid = isPageInvalid || (len(field.Errors)>0); 
		}
	}

	if isPageInvalid {
		return &ValidationError{"Form is invalid. Please update form"}
	}
	return nil
}