package services

import (
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/mappers"
	"fmt"
	"time"
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
	Util "github.com/cts3njitedu/healthful-heart/utils"
	Merge "github.com/cts3njitedu/healthful-heart/mergers"
	"strconv"
	"strings"
	"sort"
	"github.com/cts3njitedu/healthful-heart/validators"
)
type WorkoutService struct {
	locationService ILocationService
	workoutDayRepository mysqlrepo.IWorkoutDayRepository
	workoutRepository mysqlrepo.IWorkoutRepository
	pageRepository mongorepo.IPageRepository
	locationRepository mysqlrepo.ILocationRepository
	workoutTypeService IWorkoutTypeService
	groupRepository mysqlrepo.IGroupRepository
	mapper mappers.IMapper
}


func NewWorkoutService(locationService ILocationService, workoutDayRepository mysqlrepo.IWorkoutDayRepository, 
	workoutRepository mysqlrepo.IWorkoutRepository, pageRepository mongorepo.IPageRepository, locationRepository mysqlrepo.ILocationRepository, 
	workoutTypeService IWorkoutTypeService, groupRepository mysqlrepo.IGroupRepository, mapper mappers.IMapper) *WorkoutService {
	return &WorkoutService{locationService, workoutDayRepository, workoutRepository, pageRepository, locationRepository, workoutTypeService, groupRepository, mapper}
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
	categoriesAndWorkouts := serv.workoutTypeService.GetCategoriesAndWorkoutTypes();
	_, categoryCdToName := serv.workoutTypeService.GetCategories();
	for w := range heartRequest.WorkoutDays {
		workoutDay := &heartRequest.WorkoutDays[w]
		workoutDayFields, isFieldErrors := validateField(actionMap, "WORKOUTS_ACTION_PAGE.HEADER_SECTION", workoutDay.Fields)
		workoutDay.Fields = workoutDayFields;
		isError = isError || isFieldErrors
		for wd := range workoutDay.Workouts {
			workout := &workoutDay.Workouts[wd];
			// fmt.Printf("WorkoutField: %+v\n", workout.Fields)
			fillCategoryAndWorkoutType(categoriesAndWorkouts, categoryCdToName, workout.Fields)
			workoutFields, isFieldErrors := validateField(actionMap, "WORKOUTS_ACTION_PAGE.WORKOUT_SECTION", workout.Fields);
			workout.Fields = workoutFields;
			isError = isError || isFieldErrors;
			for g := range workout.Groups {
				group := &workout.Groups[g];
				groupFields, isFieldErrors := validateField(actionMap, "WORKOUTS_ACTION_PAGE.GROUP_SECTION", group.Fields);
				group.Fields = groupFields
				isError = isError || isFieldErrors
			}
		}
	}
	heartResponse := models.HeartResponse{};
	if isError {
		heartResponse.ActionType = models.WORKOUTS_ACTION_ERRORS
		heartResponse.WorkoutDays = heartRequest.WorkoutDays
	} else {
		workoutDaysCurrent, deletedIds := serv.mapper.MapWorkoutDayRequestToWorkoutDay(heartRequest, cred.UserId)
		workoutDaysOriginal := loadWorkoutDayOriginal(workoutDaysCurrent, serv)
		eventDetails := make([]models.ModEventDetail,0,10);
		workoutDaysModified := findWorkoutDaysDifferences(workoutDaysCurrent, workoutDaysOriginal, &eventDetails)
		findDeletedIds(workoutDaysOriginal, deletedIds, &eventDetails)
		fmt.Printf("Changes: %+v, Modified: %+v\n", &eventDetails, workoutDaysModified)
		for w := range workoutDaysModified {
			workoutDay := &workoutDaysModified[w]
			err := serv.workoutDayRepository.SaveWorkoutDay(workoutDay)
			if err != nil {
				addError(heartRequest.ActionType)
			}
		}
		fmt.Printf("Workout Current: %+v\n", workoutDaysModified)
		findWorkoutDaysAdded(workoutDaysModified, workoutDaysOriginal, &eventDetails)
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

func findWorkoutDaysAdded(currs [] models.WorkoutDay, origins []models.WorkoutDay, eventDetails *[]models.ModEventDetail) {
	wMap := make(map[int64]models.WorkoutDay)
	for _, wd := range origins {
		wMap[wd.Workout_Day_Id] = wd;
	}
	for _, wd := range currs {
		if val, ok := wMap[wd.Workout_Day_Id]; !ok {
			modDetail := models.ModEventDetail{};
			modDetail.Gym_Id = wd.Workout_Day_Id
			modDetail.Table_Name = "WorkoutDay"
			modDetail.Action = "ADD"
			*eventDetails = append(*eventDetails, modDetail)
			workouts := wd.Workouts;
			for _, wk := range workouts {
				modDetail := models.ModEventDetail{};
				modDetail.Gym_Id = wk.Workout_Id
				modDetail.Table_Name = "Workout"
				modDetail.Action = "ADD"
				*eventDetails = append(*eventDetails, modDetail)
				groups := wk.Groups;
				for _, g := range groups {
					modDetail := models.ModEventDetail{};
					modDetail.Gym_Id = g.Group_Id
					modDetail.Table_Name = "Group"
					modDetail.Action = "ADD"
					*eventDetails = append(*eventDetails, modDetail)
				}
			}
		} else {
			wkMap := make(map[int64]models.Workout)
			for _, wk := range val.Workouts {
				wkMap[wk.Workout_Id] = wk;
			}
			for _, wk := range wd.Workouts {
				if val, ok := wkMap[wk.Workout_Id]; !ok {
					modDetail := models.ModEventDetail{};
					modDetail.Gym_Id = wk.Workout_Id
					modDetail.Table_Name = "Workout"
					modDetail.Action = "ADD"
					*eventDetails = append(*eventDetails, modDetail)
					for _, g := range wk.Groups {
						modDetail := models.ModEventDetail{};
						modDetail.Gym_Id = g.Group_Id
						modDetail.Table_Name = "Group"
						modDetail.Action = "ADD"
						*eventDetails = append(*eventDetails, modDetail)
					}
				} else {
					gMap := make(map[int64]models.Group)
					for _, g := range val.Groups {
						gMap[g.Group_Id] = g;
					}

					for _, g := range wk.Groups {
						if _, ok := gMap[g.Group_Id]; !ok {
							modDetail := models.ModEventDetail{};
							modDetail.Gym_Id = g.Group_Id
							modDetail.Table_Name = "Group"
							modDetail.Action = "ADD"
							*eventDetails = append(*eventDetails, modDetail)
						}
					}
				}
			}
		}

	}
}
func findDeletedIds(origins [] models.WorkoutDay, deletedIds map[string][]string, eventDetails *[]models.ModEventDetail) {
	for k, ids := range deletedIds {
		for _, v := range ids {
			modDetail := models.ModEventDetail{}
			gymId, _ := strconv.ParseInt(v,10, 64)
			modDetail.Gym_Id = gymId
			modDetail.Table_Name = k;
			modDetail.Table_Column = k + "_Id"
			modDetail.Action = "DELETE"
			*eventDetails = append(*eventDetails, modDetail)
		}
	}
	
}
func findWorkoutDaysDifferences(currs [] models.WorkoutDay, origins []models.WorkoutDay, eventDetails *[]models.ModEventDetail) []models.WorkoutDay {
	wMap := make(map[int64]models.WorkoutDay)
	for _, wd := range origins {
		wMap[wd.Workout_Day_Id] = wd;
	}
	modifyWorkoutDays := make([]models.WorkoutDay, 0, len(currs))

	for _, wd := range currs {
		isModified := false;
		if val, ok := wMap[wd.Workout_Day_Id]; ok {
			if (wd.Workout_Date != val.Workout_Date) {
				modDetail := models.ModEventDetail{};
				modDetail.Gym_Id = wd.Workout_Day_Id
				modDetail.Table_Name = "WorkoutDay"
				modDetail.Table_Column = "Workout_Date"
				modDetail.Action = "MODIFY"
				oldValue := &val.Workout_Date
				modDetail.Old_Value = oldValue
				newValue := &wd.Workout_Date
				modDetail.New_Value = newValue
				*eventDetails = append(*eventDetails, modDetail)
				isModified = true;
			} else if (wd.Location_Id != val.Location_Id) {
				modDetail := models.ModEventDetail{};
				modDetail.Gym_Id = wd.Workout_Day_Id
				modDetail.Table_Name = "WorkoutDay"
				modDetail.Table_Column = "Location_Id"
				modDetail.Action = "MODIFY"
				oldValue := strconv.FormatInt(int64(val.Location_Id), 10)
				modDetail.Old_Value = &oldValue
				newValue := strconv.FormatInt(int64(wd.Location_Id), 10)
				modDetail.New_Value = &newValue
				*eventDetails = append(*eventDetails, modDetail)
				isModified = true;
			}
			modifyWorkouts := findWorkoutsDifferences(wd.Workouts, val.Workouts, eventDetails)
			if isModified || len(modifyWorkouts) > 0 {
				wd.Workouts = modifyWorkouts
				// fmt.Printf("Changed Workouts: %+v\n", wd)
				modifyWorkoutDays = append(modifyWorkoutDays, wd)
			}
			
		} else {
			modifyWorkoutDays = append(modifyWorkoutDays, wd)
		}
	}
	return modifyWorkoutDays
}

func findWorkoutsDifferences(currs [] models.Workout, origins []models.Workout, eventDetails *[]models.ModEventDetail) []models.Workout {
	wMap := make(map[int64]models.Workout)
	for _, wk := range origins {
		wMap[wk.Workout_Id] = wk;
	}
	modifyWorkouts := make([]models.Workout, 0, len(currs))
	for _, wk := range currs {
		isModified := false;
		
		if val, ok := wMap[wk.Workout_Id]; ok {
			if (wk.Workout_Type_Cd != val.Workout_Type_Cd) {
				modDetail := models.ModEventDetail{};
				modDetail.Gym_Id = wk.Workout_Id
				modDetail.Table_Name = "Workout"
				modDetail.Table_Column = "Workout_Type_Cd"
				modDetail.Action = "MODIFY"
				oldValue := &val.Workout_Type_Cd
				modDetail.Old_Value = oldValue
				newValue := &wk.Workout_Type_Cd
				modDetail.New_Value = newValue
				*eventDetails = append(*eventDetails, modDetail)
				isModified = true;
			}
			modifyGroups := findGroupDifferences(wk.Groups, val.Groups, eventDetails)
			if isModified || len(modifyGroups) > 0 {
				wk.Groups = modifyGroups
				modifyWorkouts = append(modifyWorkouts, wk)
			}
		} else {
			modifyWorkouts = append(modifyWorkouts, wk)
		}
	}
	
	return modifyWorkouts
}

func findGroupDifferences(currs [] models.Group, origins []models.Group, eventDetails *[]models.ModEventDetail) []models.Group {
	gMap := make(map[int64]models.Group)
	for _, g := range origins {
		gMap[g.Group_Id] = g;
	}

	modifyGroups := make([]models.Group, 0, len(currs))
	for _, g := range currs {
		isModified := false
		if val, ok := gMap[g.Group_Id]; ok {
			if (g.Sets != val.Sets) {
				modDetail := models.ModEventDetail{};
				modDetail.Gym_Id = g.Group_Id
				modDetail.Table_Name = "Group"
				modDetail.Table_Column = "Sets"
				modDetail.Action = "MODIFY"
				oldValue := strconv.FormatInt(int64(val.Sets), 10)
				modDetail.Old_Value = &oldValue
				newValue := strconv.FormatInt(int64(g.Sets), 10)
				modDetail.New_Value = &newValue
				*eventDetails = append(*eventDetails, modDetail)
				isModified = true;
			} else if (g.Repetitions != val.Repetitions) {
				modDetail := models.ModEventDetail{};
				modDetail.Gym_Id = g.Group_Id
				modDetail.Table_Name = "Group"
				modDetail.Table_Column = "Repetitions"
				modDetail.Action = "MODIFY"
				oldValue := strconv.FormatInt(int64(val.Repetitions), 10)
				modDetail.Old_Value = &oldValue
				newValue := strconv.FormatInt(int64(g.Repetitions), 10)
				modDetail.New_Value = &newValue
				*eventDetails = append(*eventDetails, modDetail)
				isModified = true;
			} else if (g.Weight != val.Weight) {
				modDetail := models.ModEventDetail{};
				modDetail.Gym_Id = g.Group_Id
				modDetail.Table_Name = "Group"
				modDetail.Table_Column = "Weight"
				modDetail.Action = "MODIFY"
				oldValue := strconv.FormatFloat(float64(val.Weight), 'f', -1, 32)
				modDetail.Old_Value = &oldValue
				newValue := strconv.FormatFloat(float64(g.Weight), 'f', -1, 32)
				modDetail.New_Value = &newValue
				*eventDetails = append(*eventDetails, modDetail)
				isModified = true;
			} else if (g.Duration != val.Duration) {
				modDetail := models.ModEventDetail{};
				modDetail.Gym_Id = g.Group_Id
				modDetail.Table_Name = "Group"
				modDetail.Table_Column = "Duration"
				modDetail.Action = "MODIFY"
				oldValue := strconv.FormatFloat(float64(val.Duration), 'f', -1, 32)
				modDetail.Old_Value = &oldValue
				newValue := strconv.FormatFloat(float64(g.Duration), 'f', -1, 32)
				modDetail.New_Value = &newValue
				*eventDetails = append(*eventDetails, modDetail)
				isModified = true;
			} else if (g.Variation != val.Variation) {
				modDetail := models.ModEventDetail{};
				modDetail.Gym_Id = g.Group_Id
				modDetail.Table_Name = "Group"
				modDetail.Table_Column = "Variation"
				modDetail.Action = "MODIFY"
				oldValue := &val.Variation
				modDetail.Old_Value = oldValue
				newValue := &g.Variation
				modDetail.New_Value = newValue
				*eventDetails = append(*eventDetails, modDetail)
				isModified = true;
			}
			if isModified {
				modifyGroups = append(modifyGroups, g)
			}
		} else {
			modifyGroups = append(modifyGroups, g)
		}
	}
	return modifyGroups
}
func loadWorkoutDayOriginal(workoutDaysCurrent []models.WorkoutDay, serv*WorkoutService) []models.WorkoutDay {
	workoutDaysOriginal := make([]models.WorkoutDay, 0, len(workoutDaysCurrent)+5)
	for _, workoutDay := range workoutDaysCurrent {
		workoutDayOptions := models.QueryOptions{}
		workoutDayOptions.Where = map[string]interface{} {
			"Workout_Day_Id" : workoutDay.Workout_Day_Id,
		}
		workoutDayOptions.IsEqual = true
		fmt.Printf("Workout Days: %+v\n", workoutDayOptions)
		workoutDays, _ := serv.workoutDayRepository.GetWorkoutDaysByParams(workoutDayOptions)
		
		workoutDayOriginal := models.WorkoutDay{}
		if (len(workoutDays) == 1) {
			workoutDayOriginal = workoutDays[0];
			workoutOptions := models.QueryOptions{}
			workoutOptions.Where = map[string]interface{} {
				"Workout_Day_Id" : workoutDayOriginal.Workout_Day_Id,
			}
			workoutOptions.IsEqual = true;
			workouts, _ := serv.workoutRepository.GetWorkoutByParams(workoutOptions)
			for wk := range workouts {
				wkOut := &workouts[wk]
				groupOptions := models.QueryOptions{};
				groupOptions.Where = map[string]interface{} {
					"Workout_Id" : wkOut.Workout_Id,
				}
				groupOptions.IsEqual = true;
				groups, _ := serv.groupRepository.GetGroupByParams(groupOptions)
				wkOut.Groups = groups;
			}
			workoutDayOriginal.Workouts = workouts
			workoutDaysOriginal = append(workoutDaysOriginal, workoutDayOriginal)
		}
	}
	return workoutDaysOriginal
}
func addError(actionType string) (models.HeartResponse, error) {
	heartResponse := models.HeartResponse{}
	heartResponse.ActionType = models.WORKOUTS_ACTION_ERRORS
	heartResponse.Message = actionType + ": Request Invalid"
	return heartResponse, nil
} 

func validateField(actionMap map[string]map[string]models.Field, sectionId string, fields []models.Field) ([]models.Field, bool) {
	newFields := make([]models.Field, 0, len(fields))
	fieldValidator := validators.FieldValidator{}
	var isError bool = false;
	for f := range fields {
		field := fields[f];
		newField := Util.CloneField(actionMap[sectionId][field.Name]);
		newField.Value = field.Value;
		newField.Items = field.Items
		fieldValidator.FieldValidators(&newField);
		isError = isError || len(newField.Errors) > 0;
		newFields = append(newFields, newField);
	}
	return newFields, isError
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

func fillCategoryAndWorkoutType(catWorkouts map[string]map[string]string, categories map[string]string, fields []models.Field) {
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

func fillGroupSection(groupSection models.Section, groups []models.Group) ([]models.SectionInfo) {
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