package services

import (
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
)
type GymRepositoryService struct {
	workoutDayRepository mysqlrepo.IWorkoutDayRepository
	workoutRepository mysqlrepo.IWorkoutRepository
	groupRepository mysqlrepo.IGroupRepository
	locationRepository mysqlrepo.ILocationRepository
}

func NewGymRepositoryService(workoutDayRepository mysqlrepo.IWorkoutDayRepository, workoutRepository mysqlrepo.IWorkoutRepository, 
	groupRepository mysqlrepo.IGroupRepository, locationRepository mysqlrepo.ILocationRepository) *GymRepositoryService {
	return &GymRepositoryService{workoutDayRepository, workoutRepository, groupRepository, locationRepository}
}


func (serv * GymRepositoryService) LoadWorkoutDayOriginal(workoutDaysCurrent []models.WorkoutDay) []models.WorkoutDay {
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

func (serv * GymRepositoryService) GetWorkoutDaysByParams(options models.QueryOptions) ([]models.WorkoutDay, error) {
	return serv.workoutDayRepository.GetWorkoutDaysByParams(options)
}

func (serv * GymRepositoryService) SaveWorkoutDayLocation(workDay *models.WorkoutDay) (*models.WorkoutDay, error) {
	return serv.workoutDayRepository.SaveWorkoutDayLocation(workDay)
}

func (serv * GymRepositoryService) UpdateAllWorkoutDay(workDays []models.WorkoutDay, ids map[string][]string) error  {
	return serv.workoutDayRepository.UpdateAllWorkoutDay(workDays, ids)
}

func (serv * GymRepositoryService) GetWorkoutByParams(queryOptions models.QueryOptions) ([]models.Workout, error)  {
	return serv.workoutRepository.GetWorkoutByParams(queryOptions)
}

func (serv * GymRepositoryService) GetGroupByParams(queryOptions models.QueryOptions) ([]models.Group, error) {
	return serv.groupRepository.GetGroupByParams(queryOptions)
}

func (serv * GymRepositoryService) GetLocationsQueryParams(queryOptions models.QueryOptions) ([]models.Location, error) {
	return serv.locationRepository.GetLocationsQueryParams(queryOptions)
}
