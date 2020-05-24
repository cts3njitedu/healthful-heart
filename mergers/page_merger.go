package mergers

import (

	"github.com/cts3njitedu/healthful-heart/models"
	"time"
	"fmt"
	"strconv"
	Util "github.com/cts3njitedu/healthful-heart/utils"
	"strings"
	"sort"
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

func MergeGroupToSection(section models.Section, actionType string) (models.Section) {
	for f := range section.Fields {
		var field = &section.Fields[f];
		if (field.FieldId == "EDIT" || field.FieldId == "DELETE") {
			switch actionType {
			case models.VIEW_WORKOUT_DETAILS:
				field.IsDisabled = true;
				field.IsHidden = true;
			case models.VIEW_WORKOUT_DETAILS_META_INFO: 
				field.IsDisabled = true;
				field.IsHidden = true;
			default:
				field.IsDisabled = false;
				field.IsHidden = false;
			}
		} else if (field.FieldId == "CLOSE" || field.FieldId == "SAVE") {
			switch actionType {
				case models.VIEW_WORKOUT_DETAILS: 
					field.IsDisabled = false;
					field.IsHidden = false;
				case models.VIEW_WORKOUT_DETAILS_META_INFO: 
					field.IsDisabled = true;
					field.IsHidden = true;
				default:
					field.IsDisabled = true;
					field.IsHidden = true;
			}
		}
	}
	return section;
}
func MergeWorkoutDetailsActivityToSection(section models.Section, actionType string) (models.Section) {
	for f := range section.Fields {
		var field = &section.Fields[f]
		if field.FieldId == "ADD_GROUP" {
			switch actionType {
			case models.VIEW_WORKOUT_DETAILS_META_INFO:
				field.IsDisabled = true;
				field.IsHidden = false;
			default:
				field.IsDisabled = true;
				field.IsHidden = true;
			}
		} else if (field.FieldId == "CANCEL_CHANGES") {
			switch actionType {
			case models.VIEW_WORKOUT_DETAILS_META_INFO:
				field.IsDisabled = true;
				field.IsHidden = false;
			default:
				field.IsDisabled = true;
				field.IsHidden = false;
			}
		} else if (field.FieldId == "CLOSE") {
			switch actionType {
			case models.VIEW_WORKOUT_DETAILS_META_INFO:
				field.IsDisabled = false;
				field.IsHidden = false;
			default:
				field.IsDisabled = true;
				field.IsHidden = true;
			}
		} else if (field.FieldId == "SUBMIT_CONTINUE") {
			switch actionType {
			case models.VIEW_WORKOUT_DETAILS_META_INFO:
				field.IsDisabled = true;
				field.IsHidden = false;
			default:
				field.IsDisabled = true;
				field.IsHidden = true;
			}
		} else if (field.FieldId == "SUBMIT_CLOSE") {
			switch actionType {
			case models.VIEW_WORKOUT_DETAILS_META_INFO:
				field.IsDisabled = true;
				field.IsHidden = false;
			default:
				field.IsDisabled = true;
				field.IsHidden = true;
			}
		} else if (field.FieldId == "SAVE_GROUP" || field.FieldId == "CANCEL") {
			switch actionType {
			case models.VIEW_WORKOUT_DETAILS_META_INFO:
				field.IsDisabled = true;
				field.IsHidden = true;
			default:
				field.IsDisabled = true;
				field.IsHidden = true; 
			}
		}
	}
	return section
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

func FillLocationHeaderSection(locationHeaderSection models.Section, locationSection models.Section, heartRequest models.HeartRequest) (models.SectionInfo, models.Section) {
	locationSectionInfo := models.SectionInfo{}
	newLocationSection := Util.CloneSection(locationSection)
	newLocationHeaderSection := Util.CloneSection(locationHeaderSection)
	newLocationHeaderSection.Fields = append(newLocationHeaderSection.Fields, newLocationSection.Fields...)
	cleanLocationHeaderSection := Util.CloneSection(newLocationHeaderSection)
	heartSorts := heartRequest.HeartSort;
	for f := range newLocationHeaderSection.Fields {
		field := &newLocationHeaderSection.Fields[f];
		
		if v, ok := heartSorts[field.Name]; ok {
			value := strings.ToUpper(v.SortOrder);
			field.SortOrder = &value

		}
	}
	locationMetaData := models.SectionMetaData{}
	locationMetaData.Id = heartRequest.Date;
	locationMetaData.Page = heartRequest.HeartPagination.Page
	locationSectionInfo.SectionMetaData = locationMetaData
	locationSectionInfo.Section = newLocationHeaderSection
	return locationSectionInfo, cleanLocationHeaderSection
}

func FillFilterSection(filterSection models.Section, locationSection models.Section, heartRequest models.HeartRequest) (models.SectionInfo, models.Section) {
	tableHeaders := []string {"Name", "Country", "State", "City", "ZipCode", "Location"}
	filterSectionInfo := models.SectionInfo{}
	heartFilter:= heartRequest.HeartFilter;
	newFilterSection := Util.CloneSection(filterSection);
	newLocationSection := Util.CloneSection(locationSection);
	newFilterSection.Fields = newLocationSection.Fields;
	cleanFilterSection := Util.CloneSection(newFilterSection)
	for f := range newFilterSection.Fields {
		field := &newFilterSection.Fields[f];
		if v, ok := heartFilter[field.Name]; ok {
			field.Value = v.(string)
		}
	}
	filterSectionMetaData := models.SectionMetaData{}
	filterSectionMetaData.Id = heartRequest.Date;
	filterSectionMetaData.Page = heartRequest.HeartPagination.Page
	filterSectionMetaData.TableHeaders = tableHeaders
	filterSectionInfo.SectionMetaData = filterSectionMetaData;
	filterSectionInfo.Section = newFilterSection;
	return filterSectionInfo, cleanFilterSection;
}

func FillLocationSection(locationSection models.Section, locations []models.Location, heartRequest models.HeartRequest, associatedIds map[string]interface{}) ([]models.SectionInfo) {
	fmt.Printf("Locations: %+v\n", locations)
	newSectionInfos := make([]models.SectionInfo, 0, len(locations))
	for _, loc := range locations {
		newSectionInfo := models.SectionInfo{}
		newSection := Util.CloneSection(locationSection)
		newSection = MergeLocationToSection(newSection, loc)
		newSection.IsHidden = false;
		locationMetaData := models.SectionMetaData{};
		locationMetaData.Id = strconv.FormatInt(loc.Location_Id,10)
		locationMetaData.AssociatedIds = associatedIds
		locationMetaData.Page = heartRequest.HeartPagination.Page
		newSectionInfo.SectionMetaData = locationMetaData
		newSectionInfo.Section = newSection
		newSectionInfos = append(newSectionInfos, newSectionInfo)
	}
	return newSectionInfos
}

func FillCategoryNavigationSection(navigationSection models.Section, heartRequest models.HeartRequest, categories map[string]string) (models.SectionInfo) {
	keys := make([]string, 0, len(categories))

	for k := range categories {
		keys = append(keys, k)
	}
	sort.Strings(keys);

	fmt.Printf("Sorted Categories: %+v\n", keys)
	newNavSection := Util.CloneSection(navigationSection);
	fields := make([]models.Field, 0, len(keys))
	for _, k := range keys {
		field := models.Field{};
		field.Value = k;
		field.Name = categories[k];
		fields = append(fields, field)
	}
	newNavSection.Fields = fields;
	newNavSectionInfo := models.SectionInfo{};
	newNavMetaData := models.SectionMetaData{}
	newNavSectionInfo.SectionMetaData = newNavMetaData
	newNavSectionInfo.Section = newNavSection
	return newNavSectionInfo

}

func FillNewWorkoutSection(workoutSection models.Section, heartRequest models.HeartRequest,categories map[string]string, 
			catWorkouts map[string]map[string]string, workouts map[string]models.Workout) (models.Section) {

	items := make([]models.Item, 0, len(categories));
	
	for catCode, workoutMap := range catWorkouts {
		values := make(map[string]models.Item)
		item := models.Item{};
		item.Id = catCode;
		item.Item = categories[catCode]
		for workType, workout := range workoutMap {
			if _, ok := workouts[workType]; !ok {
				values[workType] = models.Item{
					Id: workType,
					Item: workout,
				}
			}
		}
		if (len(values) > 0) {
			item.Values = values;
			items = append(items, item)
		}

	}
	newWorkoutSection := Util.CloneSection(workoutSection);
	for f := range newWorkoutSection.Fields {
		field := &newWorkoutSection.Fields[f]
		if (field.FieldId == "CATEGORY_NAME") {
			field.Items = items;
		}
	}

	return newWorkoutSection
}

func FillWorkoutSection(workoutSection models.Section, heartRequest models.HeartRequest,categories map[string]string, 
	catWorkouts map[string]map[string]string, workouts []models.Workout) ([]models.SectionInfo) {
	
	newSectionInfos := make([]models.SectionInfo, 0, len(workouts));
	for _, workout := range workouts {
		sectionInfo := models.SectionInfo{}
		sectionMetaData := models.SectionMetaData{};
		sectionMetaData.Id = strconv.FormatInt(workout.Workout_Id,10)
		sectionMetaData.VersionNb = workout.Version_Nb
		workType := workout.Workout_Type_Cd;
		var catCode string;
		var workoutName string;
		var catName string;
		for c, v := range catWorkouts {
			if wc, ok := v[workType]; ok {
				catCode = c;
				workoutName = wc
				catName = categories[c]
				break;
			}
		}

		catItem := models.Item{}
		catItem.Id = catCode;
		catItem.Item = catName;
		catItems := make([]models.Item, 0, 1);
		catItems = append(catItems, catItem)

		workItem := models.Item{}
		workItem.Id = workType
		workItem.Item = workoutName;
		workItems := make([]models.Item, 0, 1)
		workItems = append(workItems, workItem)
		newWorkoutSection := Util.CloneSection(workoutSection);
		for f := range newWorkoutSection.Fields {
			field := &newWorkoutSection.Fields[f]
			if (field.FieldId == "CATEGORY_NAME") {
				field.Value = catCode;
				field.Items = catItems;
			} else if field.FieldId == "WORKOUT_TYPE_DESC" {
				field.Value = workType
				field.Items = workItems;
			}
		}
		sectionInfo.SectionMetaData = sectionMetaData
		sectionInfo.Section = newWorkoutSection 

		newSectionInfos = append(newSectionInfos, sectionInfo)
	}
	return newSectionInfos
}

func FillCategoryAndWorkoutType(catWorkouts map[string]map[string]string, categories map[string]string, fields []models.Field) {
	var catCode string;
	var workoutType string;
	for f := range fields {
		field := &fields[f];
		if (field.Name == "categoryName") {
			catCode = field.Value;
		} else if field.Name == "workoutTypeDesc" {
			workoutType = field.Value;
		}
	}
	fmt.Printf("Category: %s, Workout %s\n", catCode, workoutType)
	for f := range fields {
		field := &fields[f];
		if (field.Name == "categoryName") {
			catItem := models.Item{}
			catItem.Id = catCode;
			catItem.Item = categories[catCode];
			catItems := make([]models.Item, 0, 1);
			catItems = append(catItems, catItem)
			field.Items = catItems
		} else if field.Name == "workoutTypeDesc" {
			var workoutTypeDesc string
			workoutTypeMap := catWorkouts[catCode];
			if workoutTypeMap != nil {
				workoutTypeDesc = workoutTypeMap[workoutType]
			}
			workItem := models.Item{}
			workItem.Id = workoutType
			workItem.Item = workoutTypeDesc;
			workItems := make([]models.Item, 0, 1)
			workItems = append(workItems, workItem)
			field.Items = workItems
		}
	}
}

func FillGroupSection(groupSection models.Section, groups []models.Group) ([]models.SectionInfo) {
	newSectionInfos := make([]models.SectionInfo, 0, len(groups))
	for _, group := range groups {
		newSection := Util.CloneSection(groupSection);
		for f := range newSection.Fields {
			field := &newSection.Fields[f];
			if (field.FieldId == "SETS") {
				field.Value = strconv.FormatInt(int64(group.Sets), 10)
				field.IsDisabled = true;
			} else if (field.FieldId == "REPETITIONS") {
				field.Value = strconv.FormatInt(int64(group.Repetitions), 10)
				field.IsDisabled = true;
			} else if (field.FieldId == "WEIGHT") {
				field.Value = strconv.FormatFloat(float64(group.Weight), 'f', -1, 32)
				field.IsDisabled = true;
			} else if (field.FieldId == "DURATION") {
				field.Value = strconv.FormatFloat(float64(group.Duration), 'f', -1, 32)
				field.IsDisabled = true;
			} else if (field.FieldId == "VARIATION") {
				field.Value = group.Variation
				field.IsDisabled = true;
			} else if (field.FieldId == "EDIT" || field.FieldId == "DELETE") {
				field.IsDisabled = false;
				field.IsHidden = false;
			} else if (field.FieldId == "CLOSE" || field.FieldId == "SAVE") {
				field.IsDisabled = true;
				field.IsHidden = true;
			}
		}
		newSectionInfo := models.SectionInfo{}
		sectionMetaData := models.SectionMetaData{};
		sectionMetaData.Id = strconv.FormatInt(group.Group_Id,10)
		sectionMetaData.VersionNb = group.Version_Nb
		newSectionInfo.SectionMetaData = sectionMetaData;
		newSectionInfo.Section = newSection
		newSectionInfos = append(newSectionInfos, newSectionInfo)
	}
	return newSectionInfos
}