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
	workoutDays, _ := serv.workoutDayRepository.GetWorkoutDays(cred.UserId)
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
	dbPage :=serv.pageRepository.GetPage("WORKOUT_DAY_PAGE");
	
	options := models.QueryOptions{};
	date, _ := time.Parse("20060102", queryParams.Date)

	dateFormat := date.Format("2006-01-02 15:04:05")

	whereClause := map[string] interface{} {
		"user_id" : cred.UserId,
		"workout_date" : dateFormat,
	}
		
	

	options.Where = whereClause;

	workouts, _ := serv.workoutDayRepository.GetWorkoutDaysByParams(options)
	fmt.Printf("workouts: %+v\n", workouts)
	headerSection := Util.FindSection("HEADER_SECTION", dbPage)
	locationSection := Util.FindSection("LOCATION_SECTION", dbPage)


	newSectionInfos := make([]models.SectionInfo, 0, len(workouts))

	newSections := make([]models.Section, 0, len(workouts))

	newSections = append(newSections, locationSection)

	newSections = append(newSections, headerSection)

	for _, workoutDay := range workouts {
		newSectionInfo := models.SectionInfo{}
		newSection := Util.CloneSection(locationSection)
		location, _ := serv.locationService.GetLocation(workoutDay.Location_Id)
		newSection = Merge.MergeLocationToSection(newSection, location)
		locationMetaData := models.SectionMetaData{};
		locationMetaData.Id = strconv.FormatInt(location.Location_Id,10)
		newSectionInfo.SectionMetaData = locationMetaData
		newSectionInfo.Section = newSection
		newSectionInfos = append(newSectionInfos, newSectionInfo)
	}
	newHeaderSection := Util.CloneSection(headerSection);
	workoutDayHeader := models.WorkoutDay{}
	dateFormat = date.Format("2006-01-02")
	workoutDayHeader.Workout_Date = dateFormat
	newHeaderSection = Merge.MergeWorkDayToSection(newHeaderSection, workoutDayHeader)
	newHeaderMetaData := models.SectionMetaData{};
	newHeaderMetaData.Id = queryParams.Date
	newHeaderSectionInfo := models.SectionInfo{}
	newHeaderSectionInfo.SectionMetaData=newHeaderMetaData
	newHeaderSectionInfo.Section = newHeaderSection
	newSectionInfos = append(newSectionInfos, newHeaderSectionInfo)
	heartResponse := models.HeartResponse{};
	heartResponse.NewSections = newSections;
	heartResponse.SectionInfos = newSectionInfos;

	return heartResponse, nil

}

func (serv * WorkoutService) GetWorkoutDaysLocationsView(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error) {
	dbPage :=serv.pageRepository.GetPage("WORKOUT_DAY_LOCATIONS_PAGE");
	locationOptions := models.QueryOptions{};
	locationWhereClause := heartRequest.HeartFilter;
	locationOptions.Where = locationWhereClause;
	filterSection := Util.FindSection("LOCATION_SECTION", dbPage)
	locationOptions.Order = Util.QueryBuildSort(heartRequest.HeartSort, filterSection)
	locationOptions.Select = [] string {"state", "city"};
	locations, _ := serv.locationRepository.GetLocationsQueryParams(locationOptions);

	fmt.Printf("Locations: %+v\n", locations)
	// dbPage :=serv.pageRepository.GetPage("WORKOUT_DAY_LOCATIONS_PAGE");

	// date, _ := time.Parse("20060102", heartRequest.Date)

	// headerSection := Util.FindSection("HEADER_SECTION", dbPage)
	// locationSection := Util.FindSection("LOCATION_SECTION", dbPage)
	// filterSection := Util.FindSection("FILTER_SECTION", dbPage)
	// activitySection := Util.FindSection("ACTIVITY_SECTION", dbPage)

	

	// newSectionInfos := make([]models.SectionInfo, 0, 5)
	// newSections := make([]models.Section, 0, 5);
	// newSections = append(newSections, locationSection)
	// newSections = append(newSections, headerSection)
	// totalLength := 0;

	// filterSectionInfo := fillFilterSection(filterSection, locationSection, heartRequest);
	// newSections = append(newSections, Util.CloneSection(filterSectionInfo.Section))

	// totalLength = totalLength + 1;

	// workouts := []models.WorkoutDay{}

	// locations := [] models.Location{}

	// if heartRequest.ActionType == models.VIEW_WORKDATE_LOCATIONS {
	// 	workouts = getWorkouts(date, cred, serv.workoutDayRepository)
	// 	ids := make([]string, 0, len(workouts))
	// 	for _, workoutDay := range workouts {
	// 		ids = append(ids, strconv.FormatInt(WorkoutDay.Location_Id, 10))
	// 	}
	// 	locations, _ = serv.locationRepository.GetByLocationIds(ids)
	// } else if heartRequest.ActionType == models.VIEW_WORKDATE_LOCATIONS {
		
	// }
	// locationOptions := models.QueryOptions{};
	// locationWhereClause := heartRequest.HeartFilter;
	// locationOptions.Where = locationWhereClause;
	// locationOptions.Order = Util.QueryBuildSort(heartRequest.HeartSort, filterSection)
	// fmt.Printf("Location Options: %+v\n", locationOptions)
	// locations, _ := locationRepository.GetLocationsQueryParams(locationOptions)

	// workoutLocation := map[int64]models.Location{};
	// for _, workoutDay := range workouts {
	// 	workoutLocation[workoutDay.Location_Id] = models.Location{}
	// }
	
	// locationSectionInfos := fillLocationSection(filterSectionInfo.Section, 
	// 	locationSection, heartRequest, workoutLocation, serv.locationRepository);
	// totalLength = totalLength + len(locationSectionInfos)	

	

	// newHeaderSection := Util.CloneSection(headerSection);
	// workoutDayHeader := models.WorkoutDay{}
	// dateFormat = date.Format("2006-01-02")
	// workoutDayHeader.Workout_Date = dateFormat
	// newHeaderSection = Merge.MergeWorkDayToSection(newHeaderSection, workoutDayHeader)
	// newHeaderMetaData := models.SectionMetaData{};
	// newHeaderMetaData.Id = heartRequest.Date
	// newHeaderSectionInfo := models.SectionInfo{}
	// newHeaderSectionInfo.SectionMetaData=newHeaderMetaData
	// newHeaderSectionInfo.Section = newHeaderSection
	// newSectionInfos = append(newSectionInfos, newHeaderSectionInfo)
	// heartResponse := models.HeartResponse{};
	// heartResponse.NewSections = newSections;
	// heartResponse.SectionInfos = newSectionInfos;
	return models.HeartResponse{}, nil;

}


func fillFilterSection(filterSection models.Section, locationSection models.Section, heartRequest models.HeartRequest) (models.SectionInfo) {
	filterSectionInfo := models.SectionInfo{}
	newFilterSection := Util.CloneSection(filterSection);
	newFilterSection.Fields = locationSection.Fields;
	filterSectionMetaData := models.SectionMetaData{}
	filterSectionMetaData.Id = heartRequest.Date;
	filterSectionInfo.SectionMetaData = filterSectionMetaData;
	filterSectionInfo.Section = newFilterSection;
	return filterSectionInfo;
}

// func fillLocationSection(filterSection models.Section, 
// 						locationSection models.Section, 
// 						heartRequest models.HeartRequest,
// 						excludeIds map[int64]models.Location) ([]models.SectionInfo) {
	
// 	fmt.Printf("Locations: %+v\n", locations)
// 	newSectionInfos := make([]models.SectionInfo, 0, len(locations))
// 	for _, loc := range locations {
// 		newSectionInfo := models.SectionInfo{}
// 		newSection := Util.CloneSection(locationSection)
// 		newSection = Merge.MergeLocationToSection(newSection, loc)
// 		if _, ok := excludeIds[loc.Location_Id] ; ok {
// 			newSection.IsHidden = true;
// 		} else {
// 			newSection.IsHidden = false;
// 		}
// 		locationMetaData := models.SectionMetaData{};
// 		locationMetaData.Id = strconv.FormatInt(loc.Location_Id,10)
// 		newSectionInfo.SectionMetaData = locationMetaData
// 		newSectionInfo.Section = newSection
// 		newSectionInfos = append(newSectionInfos, newSectionInfo)
// 	}
// 	return newSectionInfos
// }

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