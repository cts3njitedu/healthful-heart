package validators

import (

	"github.com/cts3njitedu/healthful-heart/models"
)

func validatePage(page models.Page){

	for _,section := range page.Sections {
		for _, field:= range section.Fields {
			for _, validation := range field.Validations {
				switch validation.ValidationName {
				case "MANDATORY":
					MandatoryFieldValidator(&field, &validation)
				case "REGEX":
					RegexValueValidator(&field, &validation)
				case "LENGTH":
					LengthValidator(&field, &validation)

				}
			}
		}
	}
}