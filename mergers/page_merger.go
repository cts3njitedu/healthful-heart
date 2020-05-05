package mergers

import (

	"github.com/cts3njitedu/healthful-heart/models"
	"time"
	"fmt"
	"strconv"
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

func MergeWorkDayToSection(section models.Section, workoutDay models.WorkoutDay, actionType string) (models.Section) {
	for f := range section.Fields {
		var field  = &section.Fields[f]
		if field.FieldId == "WORKOUT_DATE" {
			field.Value = workoutDay.Workout_Date
			switch actionType {
			case models.VIEW_NON_WORKOUTDATE_LOCATIONS:
				field.IsDisabled = true;
				field.IsHidden = true;
			case models.VIEW_WORKOUTS_HEADER:
				field.IsDisabled = true;
				field.IsHidden = true;
				fmt.Println(workoutDay.Workout_Date);
				date, _ := time.Parse("2006-01-02", field.Value);
				year, month, day := date.Date()
				fmt.Printf("Ball:%v, %v, %v",year, month, day)
				weekDay := date.Weekday();
				monthString := time.Month(month)
				weekDayString := weekDay.String()
				dateString :=  weekDayString + ", " + monthString.String() + " " + strconv.FormatInt(int64(day), 10) + ", " + strconv.FormatInt(int64(year), 10)
				field.Value = dateString
			default:
				field.IsDisabled = false;
				field.IsHidden = false;
			}
		} else if field.FieldId == "CANCEL" || field.FieldId == "CHANGE_DATE" {
			switch actionType {
			case models.VIEW_WORKOUTDATE_LOCATIONS:
				field.IsDisabled = true;
				field.IsHidden = false;
			default:
				field.IsDisabled = true;
				field.IsHidden = true;
			}
		} else if field.FieldId == "LOCATION" {
			switch actionType {
			case models.VIEW_WORKOUTS_HEADER:
				location := workoutDay.Location;
				if (location != models.Location{}) {
					locationString := location.Name + ", " + location.City + ", " + location.State;
					field.Value = locationString
				} 
			default:
				field.IsDisabled = true;
			}
			
		}
	} 
	return section;
}

func MergeWorkoutActivityToSection(section models.Section, actionType string) (models.Section) {
	for f := range section.Fields {
		var field  = &section.Fields[f]
		if (field.FieldId == "ADD_WORKOUT") {
			switch actionType {
			case models.VIEW_WORKOUTS_HEADER:
				field.IsDisabled = false
				field.IsHidden = false
			default:
				field.IsDisabled = true;
				field.IsHidden = true;
			}
		} else if (field.FieldId == "GO_BACK") {
			switch actionType {
			case models.VIEW_WORKOUTS_HEADER:
				field.IsDisabled = false
				field.IsHidden = false
			default:
				field.IsDisabled = true;
				field.IsHidden = true;

			}
		}
	}
	return section
}

func MergeWorkoutDayActivityToSection(section models.Section, actionType string) (models.Section) {
	for f := range section.Fields {
		var field  = &section.Fields[f]
		if (field.FieldId == "ADD_WORKOUTDATE_LOCATION") {
			switch actionType {
			case models.VIEW_NON_WORKOUTDATE_LOCATIONS:
				field.IsDisabled = true
				field.IsHidden = false
				field.IsEditable = false;
			default:
				field.IsDisabled = true
				field.IsHidden = true;
				field.IsEditable = false;
			}
		} else if field.FieldId == "VIEW_WORKOUTS" {
			switch actionType {
			case models.VIEW_WORKOUTDATE_LOCATIONS:
				field.IsDisabled = true;
				field.IsHidden = false;
				field.IsEditable = false;
			default:
				field.IsDisabled = true;
				field.IsHidden = true;
				field.IsEditable = false;
			}
		} else if field.FieldId == "DELETE_LOCATION" {
			switch actionType {
			case models.VIEW_WORKOUTDATE_LOCATIONS:
				field.IsDisabled = true;
				field.IsHidden = false;
				field.IsEditable = false;
			default:
				field.IsDisabled = true;
				field.IsHidden = true;
				field.IsEditable = false;
			}
		} else if field.FieldId == "VIEW_OTHER_LOCATIONS" {
			switch actionType {
			case models.VIEW_WORKOUTDATE_LOCATIONS:
				field.IsDisabled = false;
				field.IsHidden = false;
				field.IsEditable = false;
			default:
				field.IsDisabled = true;
				field.IsHidden = true;
				field.IsEditable = false;
			}
		} else if field.FieldId == "CANCEL" {
			field.IsDisabled = false;
			field.IsHidden = false
		} else if field.FieldId == "VIEW_WORKOUTDAY_LOCATIONS" {
			switch actionType {
			case models.VIEW_NON_WORKOUTDATE_LOCATIONS:
				field.IsDisabled = false
				field.IsHidden = false
				field.IsEditable = false;
			default:
				field.IsDisabled = true
				field.IsHidden = true;
				field.IsEditable = false;
			}
		}
	}
	return section 
}

func MergeWorkoutDayActivityToSectionLocationSelected(section models.Section, actionType string) (models.Section) {
	for f := range section.Fields {
		var field  = &section.Fields[f]
		if (field.FieldId == "ADD_WORKOUTDATE_LOCATION") {
			switch actionType {
			case models.VIEW_NON_WORKOUTDATE_LOCATIONS:
				field.IsDisabled = false
				field.IsHidden = false
				field.IsEditable = false;
			default:
				field.IsDisabled = true
				field.IsHidden = true;
				field.IsEditable = false;
			}
		} else if field.FieldId == "VIEW_WORKOUTS" {
			switch actionType {
			case models.VIEW_WORKOUTDATE_LOCATIONS:
				field.IsDisabled = false;
				field.IsHidden = false;
				field.IsEditable = false;
			default:
				field.IsDisabled = true;
				field.IsHidden = true;
				field.IsEditable = false;
			}
		} else if field.FieldId == "DELETE_LOCATION" {
			switch actionType {
			case models.VIEW_WORKOUTDATE_LOCATIONS:
				field.IsDisabled = true;
				field.IsHidden = false;
				field.IsEditable = false;
			default:
				field.IsDisabled = true;
				field.IsHidden = true;
				field.IsEditable = false;
			}
		} else if field.FieldId == "VIEW_OTHER_LOCATIONS" {
			switch actionType {
			case models.VIEW_WORKOUTDATE_LOCATIONS:
				field.IsDisabled = false;
				field.IsHidden = false;
				field.IsEditable = false;
			default:
				field.IsDisabled = true;
				field.IsHidden = true;
				field.IsEditable = false;
			}
		} else if field.FieldId == "CANCEL" {
			field.IsDisabled = false;
			field.IsHidden = false
		} else if field.FieldId == "VIEW_WORKOUTDAY_LOCATIONS" {
			switch actionType {
			case models.VIEW_NON_WORKOUTDATE_LOCATIONS:
				field.IsDisabled = false
				field.IsHidden = false
				field.IsEditable = false;
			default:
				field.IsDisabled = true
				field.IsHidden = true;
				field.IsEditable = false;
			}
		}
	}
	return section 
}