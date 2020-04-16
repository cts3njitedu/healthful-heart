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

)
type WorkoutService struct {
	locationService ILocationService
	workoutDayRepository mysqlrepo.IWorkoutDayRepository
	workoutRepository mysqlrepo.IWorkoutRepository
	pageRepository mongorepo.IPageRepository
	locationRepository mysqlrepo.ILocationRepository 
}

func NewWorkoutService(locationService ILocationService, workoutDayRepository mysqlrepo.IWorkoutDayRepository, 
	workoutRepository mysqlrepo.IWorkoutRepository, pageRepository mongorepo.IPageRepository, locationRepository mysqlrepo.ILocationRepository) *WorkoutService {
	return &WorkoutService{locationService, workoutDayRepository, workoutRepository, pageRepository, locationRepository}
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


	newSections := make([]models.Section, 0, 5);
	newSections = append(newSections, locationSection)
	newSections = append(newSections, headerSection)
	newSections = append(newSections, activitySection)
	totalLength := 0;

	filterSectionInfo := fillFilterSection(filterSection, locationSection, heartRequest);
	newSections = append(newSections, Util.CloneSection(filterSectionInfo.Section))
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
	if heartRequest.ActionType == models.VIEW_WORKOUTDATE_LOCATIONS {
		locationOptions.In = workoutQueryMap;
	} else if heartRequest.ActionType == models.VIEW_NON_WORKOUTDATE_LOCATIONS {
		locationOptions.NotIn = workoutQueryMap
	}
	locations, _ = serv.locationRepository.GetLocationsQueryParams(locationOptions)
	
	locationSectionInfos := fillLocationSection(locationSection, locations);
	totalLength = len(locationSectionInfos) + 5;
	


	newActivitySectionInfos := models.SectionInfo{};
	newActivitySectionMetaData := models.SectionMetaData{}
	newActivitySection := Util.CloneSection(activitySection)
	newActivitySectionInfos.SectionMetaData = newActivitySectionMetaData
	newActivitySectionInfos.Section = newActivitySection


	newHeaderSection := Util.CloneSection(headerSection);
	workoutDayHeader := models.WorkoutDay{}
	dateFormat = date.Format("2006-01-02")
	workoutDayHeader.Workout_Date = dateFormat
	newHeaderSection = Merge.MergeWorkDayToSection(newHeaderSection, workoutDayHeader)
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
	heartResponse := models.HeartResponse{};
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
	}
	return heartResponse, nil
	
	
}

func fillFilterSection(filterSection models.Section, locationSection models.Section, heartRequest models.HeartRequest) (models.SectionInfo) {
	filterSectionInfo := models.SectionInfo{}
	newFilterSection := Util.CloneSection(filterSection);
	newFilterSection.Fields = locationSection.Fields;
	filterSectionMetaData := models.SectionMetaData{}
	filterSectionMetaData.Id = heartRequest.Date;
	filterSectionMetaData.Page = heartRequest.HeartPagination.Page
	filterSectionInfo.SectionMetaData = filterSectionMetaData;
	filterSectionInfo.Section = newFilterSection;
	return filterSectionInfo;
}

func fillLocationSection(locationSection models.Section, locations []models.Location) ([]models.SectionInfo) {
	
	fmt.Printf("Locations: %+v\n", locations)
	newSectionInfos := make([]models.SectionInfo, 0, len(locations))
	for _, loc := range locations {
		newSectionInfo := models.SectionInfo{}
		newSection := Util.CloneSection(locationSection)
		newSection = Merge.MergeLocationToSection(newSection, loc)
		newSection.IsHidden = false;
		locationMetaData := models.SectionMetaData{};
		locationMetaData.Id = strconv.FormatInt(loc.Location_Id,10)
		newSectionInfo.SectionMetaData = locationMetaData
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