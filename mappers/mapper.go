package mappers

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"time"
	"strconv"
	"fmt"
	"strings"

)


type IMapper interface {
	MapPageToCredentials(page models.Page) models.Credentials
	MapCredentialsToUser(cred models.Credentials) models.User
	MapUserToCredentials(user models.User) models.Credentials 
	MapWorkoutDayRequestToWorkoutDay(heartRequest models.HeartRequest, userId string) ([]models.WorkoutDay, map[string][]string)
}


type Mapper struct {}
func NewMapper() *Mapper {
	return &Mapper{}
}

func (mapper *Mapper) MapPageToCredentials(page models.Page) models.Credentials {
	var cred models.Credentials
	for _,section := range page.Sections {
		for _, field := range section.Fields {
			val := *field.Value
			switch field.Name {
			case "username":
				cred.Username = val
			case "password":
				cred.Password = val
			case "confirmPassword":
				cred.ConfirmPassword = val
			case "email":
				cred.Email = val
			case "firstname":
				cred.FirstName = val
			case "lastname": 
				cred.LastName = val
			}
		}
	}
	return cred;	

}

func (mapper *Mapper) MapCredentialsToUser(cred models.Credentials) models.User {
	user:=models.User {
		Username: cred.Username,
		Password: cred.Password,
		FirstName: cred.FirstName,
		LastName: cred.LastName,
		Email: cred.Email,
	}
	return user
}

func (mapper *Mapper) MapUserToCredentials(user models.User) models.Credentials {
	credentials:=models.Credentials {
		UserId: user.User_Id,
		Username: user.Username,
		Password: user.Password,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
	}
	return credentials
}

func (mapper *Mapper) MapWorkoutDayRequestToWorkoutDay(heartRequest models.HeartRequest, userId string) ([]models.WorkoutDay, map[string][]string) {
	workoutDays := make([]models.WorkoutDay, 0, len(heartRequest.WorkoutDays))
	deletedMap := make(map[string][]string)
	user, _ := strconv.ParseInt(userId, 10, 64)
	t := time.Now()
	modTs := t.Format("2006-01-02 15:04:05")
	for _, wd := range heartRequest.WorkoutDays {
		workoutDay := models.WorkoutDay{}
		if wd.IsDeleted {
			deleteWorkoutDays := deletedMap["WorkoutDay"];
			deleteWorkoutDays = append(deleteWorkoutDays, wd.WorkoutDayId)
			deletedMap["WorkoutDay"] = deleteWorkoutDays
			continue;
		}
		workoutDayId, _ := strconv.ParseInt(wd.WorkoutDayId, 10, 64)
		if workoutDayId < 0 {
			workoutDayId = 0;
		}
		workoutDay.Workout_Day_Id = workoutDayId
		workoutDay.User_Id = user
		workoutDay.Version_Nb, _ = strconv.ParseInt(wd.Version_Nb, 10, 64)
		for f := range wd.Fields {
			field := wd.Fields[f]
			if (field.FieldId == "WORKOUT_DATE") {
				date, err := time.Parse("20060102", *field.Value)
				if err != nil {
					panic(err)
				}
				dateFormat := date.Format("2006-01-02 15:04:05")
				workoutDay.Workout_Date = dateFormat;
			} else if (field.FieldId == "LOCATION_ID") {
				locationId , err := strconv.ParseInt(*field.Value, 10, 64)
				if err != nil {
					panic(err)
				}
				workoutDay.Location_Id = locationId
				
			}
		}

		workouts := make([]models.Workout, 0, len(wd.Workouts));
		for _, wk := range wd.Workouts {
			if wk.IsDeleted {
				workoutDay.Mod_Ts = &modTs
				deleteWorkouts := deletedMap["Workout"];
				deleteWorkouts = append(deleteWorkouts, wk.WorkoutId)
				deletedMap["Workout"] = deleteWorkouts
				continue;
			}
			workout := models.Workout{};
			workoutId, _ := strconv.ParseInt(wk.WorkoutId, 10, 64)
			if workoutId < 0 {
				workoutId = 0;
			}
			workout.Workout_Id = workoutId
			workout.Workout_Day_Id = workoutDayId
			workout.Version_Nb, _ = strconv.ParseInt(wk.Version_Nb, 10, 64)
			for f := range wk.Fields {
				field := wk.Fields[f]

				if (field.FieldId == "WORKOUT_TYPE_DESC") {
					workout.Workout_Type_Id, _ = strconv.ParseInt(*field.Value, 10, 64)
				}
			}

			groups := make([]models.Group, 0, len(wk.Groups))
			fmt.Printf("Groups Mapper: %+v\n", wk.Groups)
			for _, g := range wk.Groups {
				if g.IsDeleted {
					workout.Mod_Ts = &modTs
					workoutDay.Mod_Ts = &modTs
					deleteGroups := deletedMap["Group"];
					deleteGroups = append(deleteGroups, g.GroupId)
					deletedMap["Group"] = deleteGroups
					continue;
				}
				group := models.Group{}
				groupId, _ := strconv.ParseInt(g.GroupId, 10, 64)
				if groupId < 0 {
					groupId = 0
				}
				group.Group_Id = groupId
				group.Workout_Id = workoutId
				group.Workout_Day_Id = workoutDayId
				group.Version_Nb, _ = strconv.ParseInt(g.Version_Nb, 10, 64)
				for f := range g.Fields {
					field := g.Fields[f]
					value := "";
					if field.Value != nil {
						value = *field.Value
						value = strings.TrimSpace(value)
					}
					if (field.FieldId == "SETS") {
						if field.Value == nil {
							group.Sets = nil
						} else {
							val, err := strconv.Atoi(value)
							if err != nil {
								panic(err)
							}
							group.Sets = &val;
						}
					} else if (field.FieldId == "REPETITIONS") {
						if field.Value == nil {
							group.Repetitions = nil;
						} else {
							val, err := strconv.Atoi(value)
							if err != nil {
								panic(err)
							}
							group.Repetitions = &val;
						}	
					} else if (field.FieldId == "WEIGHT") {
						if field.Value == nil || len(value) == 0 {
							group.Weight = nil;
						} else {
							weight, err := strconv.ParseFloat(value, 32)
							if err != nil {
								panic(err)
							}
							flValue := float32(weight)
							group.Weight = &flValue
						}
						
					} else if (field.FieldId == "DURATION") {
						if field.Value == nil || len(value) == 0 {
							group.Duration = nil;
						} else {
							duration, err := strconv.ParseFloat(value, 32)
							if err != nil {
								panic(err)
							}
							flValue := float32(duration)
							group.Duration = &flValue
						}
						
					} else if (field.FieldId == "VARIATION") {
						group.Variation = field.Value
					}
				}
				groups = append(groups, group)
			}
			workout.Groups = groups;
			workouts = append(workouts, workout)
		}
		workoutDay.Workouts = workouts;
		workoutDays = append(workoutDays, workoutDay)
	}
	return workoutDays, deletedMap
}