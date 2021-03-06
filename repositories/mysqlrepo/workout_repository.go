package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
	Util "github.com/cts3njitedu/healthful-heart/utils"
	"errors"
)


const SQL_QUERY_WORKOUT string = "workout_day_id = ? AND workout_type_id = ?"
type WorkoutRepository struct {
	connection connections.IMysqlConnection
	groupRepo IGroupRepository
}

func NewWorkoutRepository(connection connections.IMysqlConnection, groupRepo IGroupRepository) * WorkoutRepository {
	return &WorkoutRepository{connection, groupRepo}
}

func (repo * WorkoutRepository) DeleteWorkouts(ids map[string][]string, tx *gorm.DB) bool {
	ret := tx.Table("Workout").Where("Workout_Id IN (?)", ids["Workout"]).Or("Workout_Day_Id IN (?)", ids["WorkoutDay"]).Delete(models.Group{})
	if ret.Error != nil {
		fmt.Printf("Unable to delete Workouts %+v\n", ret.Error)
		return false;
	} else {
		fmt.Printf("Workouts Deleted: %+v\n", ret.RowsAffected)
	}
	return true;
}
func (repo *WorkoutRepository) GetWorkoutByParams(queryOptions models.QueryOptions) ([]models.Workout, error) {
	var workouts []models.Workout
	db, err := repo.connection.GetGormConnection();
	// defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	columns := map[string]models.QueryOptions {
		"workout_day_id" : models.QueryOptions{},
		"workout_id" : models.QueryOptions{},
		"workout_type_id" : models.QueryOptions{},
		"cre_ts" : models.QueryOptions{},
		"mod_ts" : models.QueryOptions{},
		"del_ts" : models.QueryOptions{},
		"version_nb" : models.QueryOptions{},
	}

	sortMap := map[string]models.QueryOptions {
		"asc" : models.QueryOptions{},
		"desc" : models.QueryOptions{},
	}
	totalQuery, values := Util.SqlQueryBuilder(queryOptions, columns, sortMap, "Workout");

	rows, err := db.Raw(totalQuery, values...).Rows()
	
	if err != nil {
		fmt.Printf("There was an error: %+v\n", err)
	} else {
		for rows.Next() {
			workout := models.Workout{};
			if err := db.ScanRows(rows, &workout); err != nil {
				fmt.Printf("Error: %+v\n", err)
			}
			workouts = append(workouts, workout)
		}
	}
	return workouts, nil;

}

func (repo * WorkoutRepository) SaveWorkout(workDay *models.Workout, tx *gorm.DB) error {
	var workoutQuery *models.Workout = &models.Workout{}
	defer func() {
		if r := recover(); r != nil {
		  tx.Rollback()
		}
	}()


	tx.Table("Workout").Where(SQL_QUERY_WORKOUT, workDay.Workout_Day_Id, workDay.Workout_Type_Id).
	First(&workoutQuery);
	
	if workDay.Workout_Id != 0 {
		if (workoutQuery.Workout_Id !=0 && (workoutQuery.Workout_Id != workDay.Workout_Id)) {
			tx.Rollback();
			fmt.Printf("Incorrect info for Workout: %+v\n", workDay.Workout_Id)
			return errors.New(fmt.Sprintf("Incorrect info for Workout: %+v\n", workDay.Workout_Id))
		}
		workoutQuery = workDay;
		workoutQuery.Version_Nb = workDay.Version_Nb
	}

	if workoutQuery.Workout_Id != 0 {
		workDay.Workout_Id = workoutQuery.Workout_Id
		ret := tx.Table("Workout").
			Where("workout_id = ? AND version_nb = ?",workDay.Workout_Id, workoutQuery.Version_Nb).
			Updates(map[string]interface{}{"mod_ts": time.Now(), "version_nb": workoutQuery.Version_Nb + 1, "workout_type_id": workoutQuery.Workout_Type_Id});
		fmt.Printf("Rows affected: %d, Workout Id: %d\n",ret.RowsAffected,workDay.Workout_Id)
		if (ret.RowsAffected == 0) {
			tx.Rollback()
			fmt.Printf("Unable to Find Workout: %+v", workDay.Workout_Id)
			return errors.New(fmt.Sprintf("Unable to Find Workout: %+v", workDay.Workout_Id))
		}	
	} else {
		t := time.Now()
		creTs := t.Format("2006-01-02 15:04:05")
		workDay.Cre_Ts = &creTs;
		workDay.Mod_Ts = nil;
		workDay.Version_Nb = 1;
		err := tx.Table("Workout").Create(&workDay).Error;
		if err != nil {
			fmt.Printf("Workout Error: %+v\n",err)
			tx.Rollback()
			return err;
		}
		fmt.Printf("Workout data: %+v\n", workDay)
		fmt.Printf("Created workout id: %d\n", workDay.Workout_Id)
	}

	for g := range workDay.Groups {
		group := &workDay.Groups[g]
		group.Workout_Id = workDay.Workout_Id
		group.Workout_Day_Id = workDay.Workout_Day_Id
		err := repo.groupRepo.SaveGroup(group, tx)
		if err != nil {
			fmt.Printf("The group error: %+v\n", group)
			tx.Rollback();
			return err;
		}
	}
	return nil;

}