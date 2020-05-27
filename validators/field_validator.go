package validators

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"regexp"
	"log"
)

type FieldValidator struct {}

func NewFieldValidator() *FieldValidator {
	return &FieldValidator{}
}

func (fieldValidator *FieldValidator) FieldValidators(field *models.Field) {
	fValidator := FieldValidator{};
	for v := range field.Validations {
		validation := &field.Validations[v]
		// log.Printf("Field Name: %+v, Validation: %+v", field.Name, validation.ValidationName)
		switch validation.ValidationName {
		case "MANDATORY":
			fValidator.MandatoryFieldValidator(field, validation)
		case "REGEX":
			fValidator.RegexValueValidator(field, validation)
		case "LENGTH":
			fValidator.LengthValidator(field, validation)
		case "DROPDOWN":
			fValidator.DropDownValidator(field, validation)

		}
	}
}

func (fieldValidator *FieldValidator) MandatoryFieldValidator(field *models.Field, v *models.Validation) {

	if field.IsMandatory {
		if field.Value == nil {
			v.IsValid = false;
		} else {
			v.IsValid = len(*field.Value)>0
		}
	} else {
		v.IsValid = true;
	}
	appendToFieldErrors(field,v)
}

func (fieldValidator *FieldValidator) RegexValueValidator(field *models.Field, v *models.Validation) {
	r,err :=regexp.Compile(field.RegexValue)
	if err!=nil {
		log.Fatal(err)
	}
	if field.Value == nil {
		if (field.MinLength != 0) {
			v.IsValid = false
		} else {
			v.IsValid = true;
		}
	} else {
		v.IsValid = r.MatchString(*field.Value);
	}
	// log.Printf("Finish Validation Regex %v", field.Name)
	appendToFieldErrors(field,v)
}

func (fieldValidator *FieldValidator) LengthValidator(field *models.Field, v *models.Validation) {
	if (field.Value == nil) {
		if (field.MinLength == 0) {
			v.IsValid = true;
		} else {
			v.IsValid = false
		}
		
	} else {
		v.IsValid = (len(*field.Value)>=field.MinLength && len(*field.Value)<=field.MaxLength)
	}
	
	appendToFieldErrors(field,v)
}

func (fieldValidator *FieldValidator) DropDownValidator(field *models.Field, v *models.Validation) {
	if (field.Items == nil || len(field.Items) == 0 || field.Value == nil) {
		v.IsValid = false;
	} else {
		v.IsValid = false;
		for i := range field.Items {
			item := field.Items[i];
			if (item.Id == *field.Value && len(item.Item) > 0) {
				v.IsValid = true;
				break;
			}
		}
	}
	appendToFieldErrors(field,v)
}

func appendToFieldErrors(field *models.Field, v *models.Validation) {
	if !v.IsValid {
		field.Errors = append(field.Errors, v.Message)
	}
}