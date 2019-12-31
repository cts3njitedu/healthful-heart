package validators

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"regexp"
	"log"
)


func MandatoryFieldValidator(field *models.Field, v *models.Validation) {

	if v.IsMandatory {
		v.IsValid = len(field.Value)>0
	} else {
		v.IsValid = true;
	}
	appendToFieldErrors(field,v)
}

func RegexValueValidator(field *models.Field, v *models.Validation) {
	r,err :=regexp.Compile(v.RegexValue)
	if err!=nil {
		log.Fatal(err)
	}
	v.IsValid = r.MatchString(field.Value);
	appendToFieldErrors(field,v)
}

func LengthValidator(field *models.Field, v *models.Validation) {
	v.IsValid = (len(field.Value)>=v.MinLength && len(field.Value)<=v.MaxLength)
	appendToFieldErrors(field,v)
}

func appendToFieldErrors(field *models.Field, v *models.Validation) {
	if !v.IsValid {
		field.Errors = append(field.Errors, v.Message)
	}
}