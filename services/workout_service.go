package services

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/mappers"
	"fmt"
	"time"
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
	Util "github.com/cts3njitedu/healthful-heart/utils"
	Merge "github.com/cts3njitedu/healthful-heart/mergers"
	"strconv"
	"github.com/cts3njitedu/healthful-heart/validators"
)
type WorkoutService struct {
	locationService ILocationService
	pageRepository mongorepo.IPageRepository
	workoutTypeService IWorkoutTypeService
	mapper mappers.IMapper
	gymRepoService IGymRepositoryService
	eventService IEventService
}


func NewWorkoutService(locationService ILocationService, pageRepository mongorepo.IPageRepository,
	workoutTypeService IWorkoutTypeService, mapper mappers.IMapper, gymRepoService IGymRepositoryService, eventService IEventService) *WorkoutService {
	return &WorkoutService{locationService, pageRepository, workoutTypeService, mapper, gymRepoService, eventService}
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
	options.IsEqual = true;
	options.WhereEqual = map[string]bool {
		"user_id": true,
	}
	options.Order = order;
	workoutDays, _ := serv.gymRepoService.GetWorkoutDaysByParams(options)
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

	filterSectionInfo, cleanFilterSection := Merge.FillFilterSection(filterSection, locationSection, heartRequest);
	newSections = append(newSections, Util.CloneSection(cleanFilterSection))

	locationHeaderSectionInfo, cleanLocationHeaderSection := Merge.FillLocationHeaderSection(locationHeaderSection, locationSection, heartRequest)
	newSections = append(newSections, Util.CloneSection(cleanLocationHeaderSection))
	locations := [] models.Location{}
	
	

	if heartRequest.ActionType == models.VIEW_WORKOUTDATE_LOCATIONS {
		options := models.QueryOptions{};
		
		options.Where = map[string]interface{} {
			"workout_date": dateFormat,
			"user_id" : cred.UserId,
		}
		for k, v := range heartRequest.HeartFilter {
			options.Where[k] = v;
		}
		
		options.WhereEqual = map[string]bool {
			"workout_date": true,
			"user_id" : true,
		}
		options.Select = []string {"workout_day_id", "location_id", "version_nb", "name", "state", "city", "country", "zipcode", "location"}
		options.IsEqual = true;
		options.Order = Util.QueryBuildSort(heartRequest.HeartSort, filterSectionInfo.Section)
		workoutLocations, err := serv.gymRepoService.GetWorkoutDaysLocationByParams(options)
		if err != nil {
			fmt.Printf("Something went horribly wrong: %+v\n", err)
		} else {
			
			for _, wkLoc := range workoutLocations {
				newLocation := models.Location{};
				newLocation.Location_Id = wkLoc.Location_Id;
				newLocation.Name = wkLoc.Location.Name;
				newLocation.State = wkLoc.Location.State;
				newLocation.City = wkLoc.Location.City;
				newLocation.Country = wkLoc.Location.Country;
				newLocation.Zipcode = wkLoc.Location.Zipcode;
				newLocation.Location = wkLoc.Location.Location
				newLocation.Workout_Day_Version_Nb = wkLoc.Version_Nb
				newLocation.Workout_Day_Id =  wkLoc.Workout_Day_Id
				locations = append(locations, newLocation)
				
			}
		}
		
	} else if heartRequest.ActionType == models.VIEW_NON_WORKOUTDATE_LOCATIONS {
		locationOptions := models.QueryOptions{};
		locationOptions.Where = heartRequest.HeartFilter;
		locationOptions.Order = Util.QueryBuildSort(heartRequest.HeartSort, filterSectionInfo.Section)
		locationOptions.Limit = heartRequest.HeartPagination.Limit
		locationOptions.Offset = heartRequest.HeartPagination.Offset
		workoutQuery := fmt.Sprintf("Select LOCATION_ID FROM WORKOUTDAY WHERE WORKOUT_DATE='%v' AND USER_ID='%v'", dateFormat,cred.UserId);
		workoutQueryMap := map[string]string {
			"location_id" : workoutQuery,
		}
		locationOptions.NotIn = workoutQueryMap
		locations, _ = serv.gymRepoService.GetLocationsQueryParams(locationOptions)
	}

	locationSectionInfos := Merge.FillLocationSection(locationSection, locations, heartRequest);
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
	workoutDay, err = serv.gymRepoService.SaveWorkoutDayLocation(workoutDay)
	
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
	selectClause := []string{"Workout_Type_Id", "Workout_Id", "Workout_Day_Id", "Version_Nb"};
	workoutOptions.Select = selectClause

	workoutOptions.Where = map[string]interface{} {
		"Workout_Day_Id" : heartRequest.WorkoutDayId,
	}
	workoutOptions.Order = map[string]string{
		"workout_id" : "asc",
	}
	workoutOptions.WhereEqual = map[string]bool{
		"Workout_Day_Id": true,
	}
	workoutOptions.IsEqual = true
	workouts, _ := serv.gymRepoService.GetWorkoutByParams(workoutOptions)
	// fmt.Printf("Workouts: %+v\n", workouts)
	typeIds := make([]int64, 0, len(workouts))
	for _, wk := range workouts {
		typeIds = append(typeIds, wk.Workout_Type_Id)
	}
	categoriesAndWorkouts := serv.workoutTypeService.GetWorkoutTypeByIds(typeIds);
	_, categoryCdToName := serv.workoutTypeService.GetCategories();
	newSectionInfos := Merge.FillWorkoutSection(workoutSection, heartRequest,categoryCdToName, categoriesAndWorkouts, workouts)

	workoutMap := make(map[int64]models.Workout)
	for _, v := range workouts {
		workoutMap[v.Workout_Type_Id] = models.Workout{}
	}

	sortedCatAndWorkouts := serv.workoutTypeService.GetSortedCategoriesAndWorkoutTypes()
	// fmt.Printf("Sorted Cats: %+v\n", sortedCatAndWorkouts)
	newWorkoutSection := Merge.FillNewWorkoutSection(workoutSection, heartRequest, sortedCatAndWorkouts, workoutMap)
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
	workoutDayOptions.WhereEqual = map[string]bool{
		"WORKOUT_DATE" : true,
		"USER_ID" : true,
		"LOCATION_ID": true,
	}
	workoutDayOptions.IsEqual = true;
	workoutDays, _ := serv.gymRepoService.GetWorkoutDaysByParams(workoutDayOptions)
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
		newHeaderMetaData.VersionNb = workoutDayHeader.Version_Nb
		newHeaderInfo.SectionMetaData = newHeaderMetaData
		newHeaderInfo.Section = newHeaderSection;
		newSectionInfos = append(newSectionInfos, newHeaderInfo)
	}

	categoryNameToCd, _ := serv.workoutTypeService.GetCategories();
	
	navSectionInfo := Merge.FillCategoryNavigationSection(navigationSection, heartRequest, categoryNameToCd)
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
	selectClause := []string{"Workout_Type_Id", "Workout_Id", "Workout_Day_Id", "Version_Nb"};
	workoutOptions.Select = selectClause

	workoutOptions.Where = map[string]interface{} {
		"Workout_Id" : heartRequest.WorkoutId,
	}
	workoutOptions.WhereEqual = map[string]bool {
		"Workout_Id": true,
	}
	workoutOptions.IsEqual = true
	workouts, _ := serv.gymRepoService.GetWorkoutByParams(workoutOptions)
	fmt.Printf("Workouts: %+v\n", workouts)
	typeIds := make([]int64, 0, len(workouts))
	for _, wk := range workouts {
		typeIds = append(typeIds, wk.Workout_Type_Id)
	}
	categoriesAndWorkouts := serv.workoutTypeService.GetWorkoutTypeByIds(typeIds);
	_, categoryCdToName := serv.workoutTypeService.GetCategories();
	
	newSectionInfos := make([]models.SectionInfo, 0, 5);

	workoutInfoSection := Merge.FillWorkoutSection(workoutSection, heartRequest,categoryCdToName, categoriesAndWorkouts, workouts)

	groupOptions := models.QueryOptions{}
	groupOptions.Where = map[string]interface{} {
		"Workout_Id" : heartRequest.WorkoutId,
	}
	groupOptions.WhereEqual = map[string]bool {
		"Workout_Id" : true,
	}
	groupOptions.IsEqual = true;
	groups, _ := serv.gymRepoService.GetGroupByParams(groupOptions);

	fmt.Printf("Groups: %+v\n", groups)
	groupInfoSections := Merge.FillGroupSection(groupSection, groups)

	newSectionInfos = append(newSectionInfos, workoutInfoSection...)

	newSectionInfos = append(newSectionInfos, groupInfoSections...)

	heartResponse := models.HeartResponse{};
	heartResponse.ActionType = heartRequest.ActionType
	heartResponse.NewSections = newSections;
	heartResponse.SectionInfos = newSectionInfos;
	return heartResponse, nil;
}
func (serv * WorkoutService) GetWorkoutDetailsMetaInfo(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error) {
	dbPage := serv.pageRepository.GetPage("WORKOUT_DETAILS_PAGE");

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

func (serv * WorkoutService) ActionWorkoutDay(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error) {
	workoutActionsPage := serv.pageRepository.GetPage("WORKOUTS_ACTION_PAGE")
	actionMap := map[string]map[string]models.Field{};
	for s := range workoutActionsPage.Sections {
		section := workoutActionsPage.Sections[s];
		fieldMap := map[string]models.Field{}
		for f := range section.Fields {
			field := section.Fields[f];
			fieldMap[field.Name] = field;
		}
		actionMap[section.Id] = fieldMap;
	}
	var isError bool = false;
	_, categoryCdToName := serv.workoutTypeService.GetCategories();
	for w := range heartRequest.WorkoutDays {
		workoutDay := &heartRequest.WorkoutDays[w]
		workoutDayFields, isFieldErrors := enrichField(actionMap, "WORKOUTS_ACTION_PAGE.HEADER_SECTION", workoutDay.Fields, workoutDay.IsDeleted)
		workoutDay.Fields = workoutDayFields;
		isError = isError || isFieldErrors
		for wd := range workoutDay.Workouts {
			workout := &workoutDay.Workouts[wd];
			fields := workout.Fields
			var wkTypeId int64;
			for f := range fields {
				field := &fields[f];
				if field.Name == "workoutTypeDesc" {
					if field.Value == nil {
						wkTypeId = 0;
					} else {
						wkTypeId, _ = strconv.ParseInt(*field.Value, 10, 64)
					}
				}
			}
			wkType := serv.workoutTypeService.GetWorkoutTypeByIds([]int64{wkTypeId});
			// if len(wkType) == 0 {
			// 	heartResponse := models.HeartResponse{};
			// 	heartResponse.SubActionType = heartRequest.SubActionType;
			// 	fmt.Printf("Something went wrong here: %v\n", wkTypeId)
			// 	heartResponse,err := addError(heartRequest.ActionType)
			// 	heartResponse.SubActionType = heartRequest.SubActionType;
			// 	return heartResponse, err
			// }
			Merge.FillCategoryAndWorkoutType(wkType[wkTypeId], categoryCdToName, workout.Fields)
			workoutFields, isFieldErrors := enrichField(actionMap, "WORKOUTS_ACTION_PAGE.WORKOUT_SECTION", workout.Fields, workout.IsDeleted);
			workout.Fields = workoutFields;
			isError = isError || isFieldErrors;
			for g := range workout.Groups {
				group := &workout.Groups[g];
				groupFields, isFieldErrors := enrichField(actionMap, "WORKOUTS_ACTION_PAGE.GROUP_SECTION", group.Fields, group.IsDeleted);
				group.Fields = groupFields
				isError = isError || isFieldErrors
			}
		}
	}
	heartResponse := models.HeartResponse{};
	heartResponse.SubActionType = heartRequest.SubActionType;
	if isError {
		heartResponse.ActionType = models.WORKOUTS_ACTION_ERRORS
		heartResponse.WorkoutDays = heartRequest.WorkoutDays
	} else {
		workoutDaysCurrent, deletedIds := serv.mapper.MapWorkoutDayRequestToWorkoutDay(heartRequest, cred.UserId)
		workoutDaysOriginal := serv.gymRepoService.LoadWorkoutDayOriginal(workoutDaysCurrent)
		eventDetails := make([]models.ModEventDetail,0,10);
		workoutDaysModified := serv.eventService.FindWorkoutDaysDifferences(workoutDaysCurrent, workoutDaysOriginal, &eventDetails)
		serv.eventService.FindDeletedIds(workoutDaysOriginal, deletedIds, &eventDetails)
		fmt.Printf("Changes: %+v, Modified: %+v, Deleted: %+v\n", &eventDetails, workoutDaysModified, deletedIds)
		err := serv.gymRepoService.UpdateAllWorkoutDay(workoutDaysModified, deletedIds)
		if err != nil {
			fmt.Printf("Something went wrong here: %+v\n", err)
			heartResponse,err = addError(heartRequest.ActionType)
			heartResponse.SubActionType = heartRequest.SubActionType;
			return heartResponse, err
		}
		fmt.Printf("Workout Current: %+v\n", workoutDaysModified)
		serv.eventService.FindWorkoutDaysAdded(workoutDaysModified, workoutDaysOriginal, &eventDetails)
		heartResponse.ActionType = models.WORKOUTS_ACTION_SUCCESS
		heartResponse.ActionInfo = buildActionResponse(eventDetails)
	}
	return heartResponse, nil;
}

func buildActionResponse(eventDetails []models.ModEventDetail) (models.ActionInfo) {
	add := make(map[string][]interface{})
	del := make(map[string][]interface{})
	modify := make(map[string][]interface{})
	for _, detail := range eventDetails {
		if (detail.Action == "ADD") {
			addIds := add[detail.Table_Name];
			add[detail.Table_Name] = append(addIds, detail.Gym_Id)
		} else if (detail.Action == "DELETE") {
			deletedIds := del[detail.Table_Name]
			del[detail.Table_Name] = append(deletedIds, detail.Gym_Id)
		} else if (detail.Action == "MODIFY") {
			modifyIds := modify[detail.Table_Name]
			modify[detail.Table_Name] = append(modifyIds, detail.Gym_Id)
		}
	}
	actionInfo := models.ActionInfo{
		Added: add,
		Deleted: del,
		Modified: modify,
	}
	return actionInfo
}

func addError(actionType string) (models.HeartResponse, error) {
	heartResponse := models.HeartResponse{}
	heartResponse.ActionType = models.WORKOUTS_ACTION_ERRORS
	heartResponse.Message = actionType + ": Request Invalid"
	return heartResponse, nil
} 

func enrichField(actionMap map[string]map[string]models.Field, sectionId string, fields []models.Field, isDeleted bool) ([]models.Field, bool) {
	newFields := make([]models.Field, 0, len(fields))
	if isDeleted {
		return newFields, false
	}
	sectionFieldMap := actionMap[sectionId];
	fieldMap := make(map[string]models.Field)
	for _, field := range fields {
		fieldMap[field.Name] = field;
	}
	fieldValidator := validators.FieldValidator{}
	var isError bool = false;
	for _, field := range sectionFieldMap {
		newField := Util.CloneField(field);
		if val, ok := fieldMap[newField.Name]; ok {
			newField.Value = val.Value;
			newField.Items = val.Items
		}
		fieldValidator.FieldValidators(&newField)
		isError = isError || len(newField.Errors) > 0;
		if _,ok := fieldMap[newField.Name]; ok {
			newFields = append(newFields, newField);
		} else {
			if len(newField.Errors) > 0 {
				newFields = append(newFields, newField);
			}
		}

	}
	return newFields, isError
}