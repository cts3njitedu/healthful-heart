package services

import (
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	"time"
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
	Util "github.com/cts3njitedu/healthful-heart/utils"
	Merge "github.com/cts3njitedu/healthful-heart/mergers"
	"strconv"
	"strings"
	"sort"

)
type WorkoutService struct {
	locationService ILocationService
	workoutDayRepository mysqlrepo.IWorkoutDayRepository
	workoutRepository mysqlrepo.IWorkoutRepository
	pageRepository mongorepo.IPageRepository
	locationRepository mysqlrepo.ILocationRepository
	workoutTypeService IWorkoutTypeService
	groupRepository mysqlrepo.IGroupRepository
}


func NewWorkoutService(locationService ILocationService, workoutDayRepository mysqlrepo.IWorkoutDayRepository, 
	workoutRepository mysqlrepo.IWorkoutRepository, pageRepository mongorepo.IPageRepository, locationRepository mysqlrepo.ILocationRepository, 
	workoutTypeService IWorkoutTypeService, groupRepository mysqlrepo.IGroupRepository) *WorkoutService {
	return &WorkoutService{locationService, workoutDayRepository, workoutRepository, pageRepository, locationRepository, workoutTypeService, groupRepository}
}

func (serv * WorkoutService) GetWorkoutDays(queryParams models.QueryParams, cred models.Credentials) ([]models.WorkoutDay, error) {
	options := models.QueryOptions{};
	whereClause := map[string] interface{} {
		"user_id" : cred.UserId,
	}
	order := map[string]string {
		"workout_date" : "asc",
	}
	options.Where = whereClause;
	options.Order = order;
	workoutDays, _ := serv.workoutDayRepository.GetWorkoutDaysByParams(options)
	fmt.Println("Retrieved workout days for user: ", cred.UserId)
	respWorkoutDays := make([]models.WorkoutDay,0 , len(workoutDays))
	for _, workoutDay := range workoutDays {
		respWorkoutDay := models.WorkoutDay{}
		respWorkoutDay.Workout_Day_Id = workoutDay.Workout_Day_Id
		respWorkoutDay.User_Id = workoutDay.User_Id
		workoutTime, _ := time.Parse("2006-01-02 15:04:05", workoutDay.Workout_Date)
		respWorkoutDay.Year = workoutTime.Year()
		month := workoutTime.Month() 
		respWorkoutDay.Month = month.String()
		respWorkoutDay.MonthId = int(workoutTime.Month())
		respWorkoutDay.Day = workoutTime.Day()
		lastDayTime := time.Date(workoutTime.Year(), workoutTime.Month() + 1,0, 0, 0, 0, 0, time.UTC)
		respWorkoutDay.NumberOfDays = lastDayTime.Day()
		location, _ := serv.locationService.GetLocation(workoutDay.Location_Id)
		respWorkoutDay.Location = location
		respWorkoutDays = append(respWorkoutDays, respWorkoutDay)
	}
	return respWorkoutDays, nil
}

func (serv * WorkoutService) GetWorkoutDaysPage(queryParams models.QueryParams, cred models.Credentials) (models.HeartResponse, error) {
	return models.HeartResponse{}, nil

}

func (serv * WorkoutService) GetWorkoutDaysLocationsView(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error) {
	dbPage :=serv.pageRepository.GetPage("WORKOUT_DAY_LOCATIONS_PAGE");
	date, _ := time.Parse("20060102", heartRequest.Date)
	dateFormat := date.Format("2006-01-02 15:04:05")


	headerSection := Util.FindSection("HEADER_SECTION", dbPage)
	locationSection := Util.FindSection("LOCATION_SECTION", dbPage)
	filterSection := Util.FindSection("FILTER_SECTION", dbPage)
	activitySection := Util.FindSection("ACTIVITY_SECTION", dbPage)
	locationHeaderSection := Util.FindSection("LOCATION_HEADER_SECTION", dbPage)


	newSections := make([]models.Section, 0, 5);
	newSections = append(newSections, locationSection)
	newSections = append(newSections, headerSection)
	
	totalLength := 0;

	filterSectionInfo, cleanFilterSection := fillFilterSection(filterSection, locationSection, heartRequest);
	newSections = append(newSections, Util.CloneSection(cleanFilterSection))

	locationHeaderSectionInfo, cleanLocationHeaderSection := fillLocationHeaderSection(locationHeaderSection, locationSection, heartRequest)
	newSections = append(newSections, Util.CloneSection(cleanLocationHeaderSection))
	locations := [] models.Location{}
	locationOptions := models.QueryOptions{};
	locationOptions.Where = heartRequest.HeartFilter;
	locationOptions.Order = Util.QueryBuildSort(heartRequest.HeartSort, filterSectionInfo.Section)
	locationOptions.Limit = heartRequest.HeartPagination.Limit
	locationOptions.Offset = heartRequest.HeartPagination.Offset
	workoutQuery := fmt.Sprintf("Select LOCATION_ID FROM WORKOUTDAY WHERE WORKOUT_DATE='%v' AND USER_ID='%v'", dateFormat,cred.UserId);
	workoutQueryMap := map[string]string {
		"location_id" : workoutQuery,
	}
	
	associatedIds := make(map[string]interface{});
	// workoutDayOptions := models.QueryOptions{}
	// workoutDayOptions.Where = map[string]interface{} {
	// 	"WORKOUT_DATE" : dateFormat,
	// 	"USER_ID" : cred.UserId,
	// }
	// workoutDayOptions.Select = []string{"Workout_Day_Id", "Location_Id"}
	// workoutDayOptions.IsEqual = true;
	// workoutDays, _ := serv.workoutDayRepository.GetWorkoutDaysByParams(workoutDayOptions)
	// workoutDayMap := make(map[string]models.WorkoutDay);
	if heartRequest.ActionType == models.VIEW_WORKOUTDATE_LOCATIONS {
		locationOptions.In = workoutQueryMap;
	} else if heartRequest.ActionType == models.VIEW_NON_WORKOUTDATE_LOCATIONS {
		locationOptions.NotIn = workoutQueryMap
	}
	locations, _ = serv.locationRepository.GetLocationsQueryParams(locationOptions)
	
	locationSectionInfos := fillLocationSection(locationSection, locations, heartRequest,associatedIds);
	totalLength = len(locationSectionInfos) + 5;
	


	newActivitySectionInfos := models.SectionInfo{};
	newActivitySectionMetaData := models.SectionMetaData{}
	newActivitySection := Util.CloneSection(activitySection)
	newActivitySection = Merge.MergeWorkoutDayActivityToSection(newActivitySection, heartRequest.ActionType)
	newActivitySectionInfos.SectionMetaData = newActivitySectionMetaData
	newActivitySectionInfos.Section = newActivitySection

	locationSelectedSection := Util.CloneSection(activitySection);
	locationSelectedSection = Merge.MergeWorkoutDayActivityToSectionLocationSelected(locationSelectedSection, heartRequest.ActionType)
	locationSelectedSection.SectionId = models.LOCATION_SELECTED
	newActivitySectionStore := Util.CloneSection(newActivitySection);
	newActivitySectionStore.SectionId = models.DEFAULT_ACTIVITY
	newSections = append(newSections, locationSelectedSection)
	newSections = append(newSections, newActivitySectionStore)


	newHeaderSection := Util.CloneSection(headerSection);
	workoutDayHeader := models.WorkoutDay{}
	dateFormat = date.Format("2006-01-02")
	workoutDayHeader.Workout_Date = dateFormat
	newHeaderSection = Merge.MergeWorkDayToSection(newHeaderSection, workoutDayHeader, heartRequest.ActionType)
	newHeaderMetaData := models.SectionMetaData{};
	newHeaderMetaData.Id = heartRequest.Date
	newHeaderSectionInfo := models.SectionInfo{}
	newHeaderSectionInfo.SectionMetaData=newHeaderMetaData
	newHeaderSectionInfo.Section = newHeaderSection

	newSectionInfos := make([]models.SectionInfo, 0, totalLength)
	newSectionInfos = append(newSectionInfos, newHeaderSectionInfo)
	newSectionInfos = append(newSectionInfos, filterSectionInfo)
	newSectionInfos = append(newSectionInfos, newActivitySectionInfos)
	newSectionInfos = append(newSectionInfos, locationSectionInfos...)
	newSectionInfos = append(newSectionInfos, locationHeaderSectionInfo)
	heartResponse := models.HeartResponse{};
	heartResponse.ActionType = heartRequest.ActionType
	heartResponse.NewSections = newSections;
	heartResponse.SectionInfos = newSectionInfos;
	return heartResponse, nil;

}

func (serv * WorkoutService) AddWorkoutDateLocation(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error) {
	date, _ := time.Parse("20060102", heartRequest.Date)
	dateFormat := date.Format("2006-01-02 15:04:05")
	sectionInfo := heartRequest.SectionInfo;
	location_id := sectionInfo.SectionMetaData.Id;
	heartResponse := models.HeartResponse{}
	var workoutDay *models.WorkoutDay = &models.WorkoutDay{};
	workoutDay.Workout_Date = dateFormat;
	workoutTime, _ := time.Parse("2006-01-02 15:04:05", workoutDay.Workout_Date)
	workoutDay.Year = workoutTime.Year()
	month := workoutTime.Month() 
	workoutDay.Month = month.String()
	workoutDay.MonthId = int(workoutTime.Month())
	workoutDay.Day = workoutTime.Day()
	lastDayTime := time.Date(workoutTime.Year(), workoutTime.Month() + 1,0, 0, 0, 0, 0, time.UTC)
	workoutDay.NumberOfDays = lastDayTime.Day()
	userId, err := strconv.ParseInt(cred.UserId, 10, 64)
	workoutDay.User_Id = userId
	id, err := strconv.ParseInt(location_id, 10, 64)
	if err != nil {
		heartResponse.Message = "Action Failed"
		heartResponse.IsSuccessful = false
		return heartResponse, nil;
	}
	workoutDay.Location_Id = id;
	workoutDay, err = serv.workoutDayRepository.SaveWorkoutDayLocation(workoutDay)
	
	if err != nil {
		heartResponse.Message = "Action Failed"
		heartResponse.IsSuccessful = false
	} else {
		heartResponse.Message = "Action Successful"
		heartResponse.IsSuccessful = true
		heartResponse.Data = workoutDay
		heartResponse.ActionType = models.VIEW_WORKOUTDATE_LOCATIONS
	}
	return heartResponse, nil
	
	
}

func (serv *WorkoutService) GetWorkouts(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error) {
	dbPage :=serv.pageRepository.GetPage("WORKOUTS_PAGE");
	workoutSection := Util.FindSection("WORKOUT_SECTION", dbPage);

	workoutOptions := models.QueryOptions{}
	selectClause := []string{"Workout_Type_Cd", "Workout_Id", "Workout_Day_Id", "Version_Nb"};
	workoutOptions.Select = selectClause

	workoutOptions.Where = map[string]interface{} {
		"Workout_Day_Id" : heartRequest.WorkoutDayId,
	}
	workoutOptions.Order = map[string]string{
		"workout_id" : "asc",
	}
	workoutOptions.IsEqual = true
	workouts, _ := serv.workoutRepository.GetWorkoutByParams(workoutOptions)
	fmt.Printf("Workouts: %+v\n", workouts)
	categoriesAndWorkouts := serv.workoutTypeService.GetCategoriesAndWorkoutTypes();
	_, categoryCdToName := serv.workoutTypeService.GetCategories();
	newSectionInfos := fillWorkoutSection(workoutSection, heartRequest,categoryCdToName, categoriesAndWorkouts, workouts)

	workoutMap := make(map[string]models.Workout)
	for _, v := range workouts {
		workoutMap[v.Workout_Type_Cd] = models.Workout{}
	}

	newWorkoutSection := fillNewWorkoutSection(workoutSection, heartRequest, categoryCdToName, categoriesAndWorkouts, workoutMap)
	newSections := make([]models.Section, 0, 5);
	newSections = append(newSections, newWorkoutSection)

	heartResponse := models.HeartResponse{};
	heartResponse.ActionType = heartRequest.ActionType
	heartResponse.SectionInfos = newSectionInfos;
	heartResponse.NewSections = newSections
	return heartResponse, nil
}
func (serv * WorkoutService) GetWorkoutPageHeader(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error) {
	dbPage :=serv.pageRepository.GetPage("WORKOUTS_PAGE");
	date, _ := time.Parse("20060102", heartRequest.Date)
	dateFormat := date.Format("2006-01-02 15:04:05")


	headerSection := Util.FindSection("HEADER_SECTION", dbPage)
	navigationSection := Util.FindSection("NAVIGATION_SECTION", dbPage)
	activitySection := Util.FindSection("ACTIVITY_SECTION", dbPage)

	newSections := make([]models.Section, 0, 5);
	
	workoutDayOptions := models.QueryOptions{}
	workoutDayOptions.Where = map[string]interface{} {
		"WORKOUT_DATE" : dateFormat,
		"USER_ID" : cred.UserId,
		"LOCATION_ID": heartRequest.LocationId,
	}
	workoutDayOptions.IsEqual = true;
	workoutDays, _ := serv.workoutDayRepository.GetWorkoutDaysByParams(workoutDayOptions)
	fmt.Printf("Workout Days: %+v\n", workoutDays)
	newSectionInfos := make([]models.SectionInfo, 0, 5);
	workoutDayHeader := models.WorkoutDay{};
	if len(workoutDays) <= 1 {
		if len(workoutDays) == 1 {
			workoutDayHeader = workoutDays[0];
			location, _ := serv.locationService.GetLocation(workoutDayHeader.Location_Id)
			workoutDayHeader.Location = location
		}
		dateFormat = date.Format("2006-01-02")
		workoutDayHeader.Workout_Date = dateFormat
		newHeaderSection := Util.CloneSection(headerSection);
		newHeaderSection = Merge.MergeWorkDayToSection(newHeaderSection, workoutDayHeader, heartRequest.ActionType)
		newSections = append(newSections, Util.CloneSection(newHeaderSection))
		newHeaderInfo := models.SectionInfo{};
		newHeaderMetaData := models.SectionMetaData{};
		workoutId := strconv.FormatInt(workoutDayHeader.Workout_Day_Id,10)
		newHeaderMetaData.Id = workoutId;
		newHeaderInfo.SectionMetaData = newHeaderMetaData
		newHeaderInfo.Section = newHeaderSection;
		newSectionInfos = append(newSectionInfos, newHeaderInfo)
	}

	categoryNameToCd, _ := serv.workoutTypeService.GetCategories();
	
	navSectionInfo := fillCategoryNavigationSection(navigationSection, heartRequest, categoryNameToCd)
	newSections = append(newSections, Util.CloneSection(navSectionInfo.Section))
	newSectionInfos = append(newSectionInfos, navSectionInfo)
	
	

	newActivitySection := Util.CloneSection(activitySection);
	newActivitySection = Merge.MergeWorkoutActivityToSection(newActivitySection, heartRequest.ActionType);
	newActivitySectionInfo := models.SectionInfo{}
	newActivitySectionInfo.SectionMetaData = models.SectionMetaData{};
	newActivitySectionInfo.Section = newActivitySection;
	newSections = append(newSections, Util.CloneSection(newActivitySection))
	newSectionInfos = append(newSectionInfos, newActivitySectionInfo)
	heartResponse := models.HeartResponse{};
	heartResponse.ActionType = heartRequest.ActionType
	heartResponse.NewSections = newSections;
	heartResponse.SectionInfos = newSectionInfos;
	return heartResponse, nil
}

func (serv * WorkoutService) GetWorkoutDetails(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error) {
	dbPage :=serv.pageRepository.GetPage("WORKOUT_DETAILS_PAGE");

	workoutSection := Util.FindSection("WORKOUT_SECTION", dbPage)
	groupSection := Util.FindSection("GROUP_SECTION", dbPage)

	newSections := make([]models.Section, 0, 5);

	workoutOptions := models.QueryOptions{}
	selectClause := []string{"Workout_Type_Cd", "Workout_Id", "Workout_Day_Id", "Version_Nb"};
	workoutOptions.Select = selectClause

	workoutOptions.Where = map[string]interface{} {
		"Workout_Id" : heartRequest.WorkoutId,
	}
	workoutOptions.IsEqual = true
	workouts, _ := serv.workoutRepository.GetWorkoutByParams(workoutOptions)
	fmt.Printf("Workouts: %+v\n", workouts)
	categoriesAndWorkouts := serv.workoutTypeService.GetCategoriesAndWorkoutTypes();
	_, categoryCdToName := serv.workoutTypeService.GetCategories();
	
	newSectionInfos := make([]models.SectionInfo, 0, 5);

	workoutInfoSection := fillWorkoutSection(workoutSection, heartRequest,categoryCdToName, categoriesAndWorkouts, workouts)

	groupOptions := models.QueryOptions{}
	groupOptions.Where = map[string]interface{} {
		"Workout_Id" : heartRequest.WorkoutId,
	}
	groupOptions.IsEqual = true;
	groups, _ := serv.groupRepository.GetGroupByParams(groupOptions);

	fmt.Printf("Groups: %+v\n", groups)
	groupInfoSections := fillGroupSection(groupSection, groups)

	newSectionInfos = append(newSectionInfos, workoutInfoSection...)

	newSectionInfos = append(newSectionInfos, groupInfoSections...)

	heartResponse := models.HeartResponse{};
	heartResponse.ActionType = heartRequest.ActionType
	heartResponse.NewSections = newSections;
	heartResponse.SectionInfos = newSectionInfos;
	return heartResponse, nil;
}
func (serv * WorkoutService) GetWorkoutDetailsMetaInfo(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error) {
	dbPage :=serv.pageRepository.GetPage("WORKOUT_DETAILS_PAGE");

	groupSection := Util.FindSection("GROUP_SECTION", dbPage)
	activitySection := Util.FindSection("ACTIVITY_SECTION", dbPage)
	

	newSections := make([]models.Section, 0, 5);

	newSectionInfos := make([]models.SectionInfo, 0, 5);

	

	newActivitySection := Util.CloneSection(activitySection);
	newActivitySection = Merge.MergeWorkoutDetailsActivityToSection(newActivitySection, heartRequest.ActionType);
	newActivitySectionInfo := models.SectionInfo{}
	newActivitySectionInfo.SectionMetaData = models.SectionMetaData{};
	newActivitySectionInfo.Section = newActivitySection;
	newSections = append(newSections, Util.CloneSection(newActivitySection))

	newGroupSection := Util.CloneSection(groupSection)
	newGroupSection = Merge.MergeGroupToSection(newGroupSection, heartRequest.ActionType);

	newSections = append(newSections, newGroupSection);
	newSectionInfos = append(newSectionInfos, newActivitySectionInfo)

	heartResponse := models.HeartResponse{};
	heartResponse.ActionType = heartRequest.ActionType
	heartResponse.NewSections = newSections;
	heartResponse.SectionInfos = newSectionInfos;
	return heartResponse, nil
}

func fillLocationHeaderSection(locationHeaderSection models.Section, locationSection models.Section, heartRequest models.HeartRequest) (models.SectionInfo, models.Section) {
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

func fillFilterSection(filterSection models.Section, locationSection models.Section, heartRequest models.HeartRequest) (models.SectionInfo, models.Section) {
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

func fillLocationSection(locationSection models.Section, locations []models.Location, heartRequest models.HeartRequest, associatedIds map[string]interface{}) ([]models.SectionInfo) {
	fmt.Printf("Locations: %+v\n", locations)
	newSectionInfos := make([]models.SectionInfo, 0, len(locations))
	for _, loc := range locations {
		newSectionInfo := models.SectionInfo{}
		newSection := Util.CloneSection(locationSection)
		newSection = Merge.MergeLocationToSection(newSection, loc)
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

func fillCategoryNavigationSection(navigationSection models.Section, heartRequest models.HeartRequest, categories map[string]string) (models.SectionInfo) {
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

func fillNewWorkoutSection(workoutSection models.Section, heartRequest models.HeartRequest,categories map[string]string, 
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

func fillWorkoutSection(workoutSection models.Section, heartRequest models.HeartRequest,categories map[string]string, 
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

func fillGroupSection(groupSection models.Section, groups []models.Group) ([]models.SectionInfo) {
	newSectionInfos := make([]models.SectionInfo, 0, len(groups))
	for _, group := range groups {
		newSection := Util.CloneSection(groupSection);
		for f := range newSection.Fields {
			field := &newSection.Fields[f];
			if (field.FieldId == "SETS") {
				field.Value = strconv.FormatInt(int64(group.Sets), 10)
			} else if (field.FieldId == "REPETITIONS") {
				field.Value = strconv.FormatInt(int64(group.Repetitions), 10)
			} else if (field.FieldId == "WEIGHT") {
				field.Value = strconv.FormatFloat(float64(group.Weight), 'f', 2, 32)
			} else if (field.FieldId == "DURATION") {
				field.Value = strconv.FormatFloat(float64(group.Duration), 'f', 2, 32)
			} else if (field.FieldId == "VARIATION") {
				field.Value = group.Variation
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

func getWorkouts(date time.Time, cred models.Credentials, workoutDayRepository mysqlrepo.IWorkoutDayRepository) ([]models.WorkoutDay) {
	options := models.QueryOptions{};

	dateFormat := date.Format("2006-01-02 15:04:05")
	fmt.Printf("Credentials: %+v\n",cred)
	whereClause := map[string] interface{} {
		"user_id" : cred.UserId,
		"workout_date" : dateFormat,
	}

	selectClause := []string{"location_id"}

	options.Where = whereClause;
	options.Select = selectClause

	workouts, _ := workoutDayRepository.GetWorkoutDaysSpecifyColumns(options)
	return workouts;
}