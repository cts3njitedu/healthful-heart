package mergers

import (

	"github.com/cts3njitedu/healthful-heart/models"
)


type PageMerger struct {}



func NewPageMerger() *PageMerger {
	return &PageMerger{}
}
func (merge *PageMerger) MergeRequestPageToPage(requestPage *models.Page, page models.Page) {
	fieldMap:=make(map[string]models.Field)
	for _, section := range page.Sections {
		for _, field := range section.Fields {
			fieldMap[field.Id] = field

		}
	}

	for s := range requestPage.Sections {
		var section = &requestPage.Sections[s];
		for f := range section.Fields {
			var field  = &section.Fields[f]
			dbField:=fieldMap[field.Id];
			field.IsMandatory = dbField.IsMandatory
			field.RegexValue = dbField.RegexValue
			field.MinLength = dbField.MinLength
			field.MaxLength = dbField.MaxLength
			field.Name = dbField.Name
			field.Validations = dbField.Validations
		}
	}


}

func MergeLocationToSection(section models.Section, location models.Location) (models.Section) {
	for f := range section.Fields {
		var field  = &section.Fields[f]
		if field.FieldId == "LOCATION" {
			field.Value = location.Location;
		} else if field.FieldId == "ZIPCODE" {
			field.Value = location.Zipcode
		} else if field.FieldId == "STATE" {
			field.Value = location.State
		} else if field.FieldId == "CITY" {
			field.Value = location.City
		} else if field.FieldId == "COUNTRY" {
			field.Value = location.Country
		} else if field.FieldId == "NAME" {
			field.Value = location.Name
		}
	}	
	return section
}

func MergeWorkDayToSection(section models.Section, workoutDay models.WorkoutDay) (models.Section) {
	for f := range section.Fields {
		var field  = &section.Fields[f]
		if (field.FieldId == "WORKOUT_DATE") {
			field.Value = workoutDay.Workout_Date
		}
	} 
	return section;
}