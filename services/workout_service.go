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
}

func NewWorkoutService(locationService ILocationService, workoutDayRepository mysqlrepo.IWorkoutDayRepository, workoutRepository mysqlrepo.IWorkoutRepository, pageRepository mongorepo.IPageRepository) *WorkoutService {
	return &WorkoutService{locationService, workoutDayRepository, workoutRepository, pageRepository}
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


	locations, _ := serv.locationService.GetLocations()

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
	newHeaderSection = Merge.MergeWorkDayAndLocationsToSection(newHeaderSection, workoutDayHeader, locations)
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

