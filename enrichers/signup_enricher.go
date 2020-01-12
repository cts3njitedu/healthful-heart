package enrichers

import (
	"github.com/cts3njitedu/healthful-heart/models"
)

type SignupEnrich struct {}


func NewSignupEnrich() *SignupEnrich {
	return &SignupEnrich{}
}


func (enrich *SignupEnrich) Enrich(page *models.Page, pageTemplate models.Page) {

	if page.PageId != "SIGNUP" && page.PageId != "LOGIN" {
		return
	}
	fieldMap:=make(map[string]models.Field)
	for _, section := range pageTemplate.Sections {
		for _, field := range section.Fields {
			fieldMap[field.Id] = field

		}
	}
	

	for s := range page.Sections {
		var section = &page.Sections[s];
		for f := range section.Fields {
			var field  = &section.Fields[f]
			var mandatory bool
			var regexValue string
			var minLength int
			var maxLength int
			dbField:=fieldMap[field.Id];
			if page.PageId == "SIGNUP" {
				mandatory = true;
				regexValue = dbField.RegexValue
				minLength = dbField.MinLength
				maxLength = dbField.MaxLength
			} else {
				mandatory = dbField.IsMandatory
				if dbField.Name == "username" || dbField.Name == "password" {
					regexValue = dbField.RegexValue
					minLength = dbField.MinLength
					maxLength = dbField.MaxLength
				} else {
					regexValue = ""
					minLength = 0
					maxLength = 0
				}
			}
			field.IsMandatory = mandatory
			field.RegexValue = regexValue
			field.MinLength = minLength
			field.MaxLength = maxLength
			field.Name = dbField.Name
			field.Validations = dbField.Validations
		}
	}
	
}