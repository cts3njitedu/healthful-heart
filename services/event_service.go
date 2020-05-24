package services

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"strconv"
)
type EventService struct {}

func NewEventService() *EventService {
	return &EventService{}
}

func (serv * EventService) FindWorkoutDaysAdded(currs [] models.WorkoutDay, origins []models.WorkoutDay, eventDetails *[]models.ModEventDetail) {
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

func (serv * EventService) FindDeletedIds(origins [] models.WorkoutDay, deletedIds map[string][]string, eventDetails *[]models.ModEventDetail) {
	for k, ids := range deletedIds {
		for _, v := range ids {
			modDetail := models.ModEventDetail{}
			gymId, _ := strconv.ParseInt(v,10, 64)
			modDetail.Gym_Id = gymId
			modDetail.Table_Name = k;
			modDetail.Action = "DELETE"
			*eventDetails = append(*eventDetails, modDetail)
		}
	}
	
}
func (serv * EventService) FindWorkoutDaysDifferences(currs [] models.WorkoutDay, origins []models.WorkoutDay, eventDetails *[]models.ModEventDetail) []models.WorkoutDay {
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
			modifyWorkouts := serv.FindWorkoutsDifferences(wd.Workouts, val.Workouts, eventDetails)
			if isModified || len(modifyWorkouts) > 0 || wd.Mod_Ts != nil {
				wd.Workouts = modifyWorkouts
				// fmt.Printf("Changed Workouts: %+v\n", wd)
				if wd.Mod_Ts != nil || len(modifyWorkouts) > 0 {
					modDetail := models.ModEventDetail{};
					modDetail.Gym_Id = wd.Workout_Day_Id
					modDetail.Table_Name = "WorkoutDay"
					if wd.Mod_Ts != nil {
						modDetail.Table_Column = "Workout Deleted"
					} else {
						modDetail.Table_Column = "Workout Added"
					}
					modDetail.Action = "MODIFY"
					*eventDetails = append(*eventDetails, modDetail)

				}
				modifyWorkoutDays = append(modifyWorkoutDays, wd)
			}
			
		} else {
			modifyWorkoutDays = append(modifyWorkoutDays, wd)
		}
	}
	return modifyWorkoutDays
}

func (serv * EventService) FindWorkoutsDifferences(currs [] models.Workout, origins []models.Workout, eventDetails *[]models.ModEventDetail) []models.Workout {
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
			modifyGroups := serv.FindGroupDifferences(wk.Groups, val.Groups, eventDetails)
			if isModified || len(modifyGroups) > 0 || wk.Mod_Ts != nil {
				wk.Groups = modifyGroups
				if wk.Mod_Ts != nil || len(modifyGroups) > 0 {
					modDetail := models.ModEventDetail{};
					modDetail.Gym_Id = wk.Workout_Id
					modDetail.Table_Name = "Workout"
					if wk.Mod_Ts != nil {
						modDetail.Table_Column = "Group Deleted"
					} else {
						modDetail.Table_Column = "Group Added"
					}
					modDetail.Action = "MODIFY"
					*eventDetails = append(*eventDetails, modDetail)

				}
				modifyWorkouts = append(modifyWorkouts, wk)
			}
		} else {
			modifyWorkouts = append(modifyWorkouts, wk)
		}
	}
	
	return modifyWorkouts
}

func (serv * EventService) FindGroupDifferences(currs [] models.Group, origins []models.Group, eventDetails *[]models.ModEventDetail) []models.Group {
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