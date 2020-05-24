package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
	Util "github.com/cts3njitedu/healthful-heart/utils"
)


const SQL_QUERY_WORKOUT string = "workout_day_id = ? AND workout_type_cd = ?"
type WorkoutRepository struct {
	connection connections.IMysqlConnection
	groupRepo IGroupRepository
}

func NewWorkoutRepository(connection connections.IMysqlConnection, groupRepo IGroupRepository) * WorkoutRepository {
	return &WorkoutRepository{connection, groupRepo}
}

func (repo *WorkoutRepository) GetWorkoutByParams(queryOptions models.QueryOptions) ([]models.Workout, error) {
	var workouts []models.Workout
	db, err := repo.connection.GetGormConnection();
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	columns := map[string]models.QueryOptions {
		"workout_day_id" : models.QueryOptions{},
		"workout_id" : models.QueryOptions{},
		"workout_type_cd" : models.QueryOptions{},
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

	if (workDay.Workout_Id == 0) {
		tx.Table("Workout").Where(SQL_QUERY_WORKOUT, workDay.Workout_Day_Id, workDay.Workout_Type_Cd).
		First(&workoutQuery);
	} else {
		
		workoutQuery = workDay
		fmt.Printf("Workout is Here: %+v\n", workoutQuery)
	}

	if workoutQuery.Workout_Id != 0 {
		workDay.Workout_Id = workoutQuery.Workout_Id
		ret := tx.Table("Workout").
			Where("workout_id = ? AND version_nb = ?",workDay.Workout_Id, workoutQuery.Version_Nb).
			Updates(map[string]interface{}{"mod_ts": time.Now(), "version_nb": workoutQuery.Version_Nb + 1, "workout_type_cd": workoutQuery.Workout_Type_Cd});
		fmt.Printf("Rows affected: %d, Workout Id: %d\n",ret.RowsAffected,workDay.Workout_Id)	
	} else {
		t := time.Now()
		creTs := t.Format("2006-01-02 15:04:05")
		workDay.Cre_Ts = &creTs;
		workDay.Version_Nb = 1;
		err := tx.Table("Workout").Create(&workDay).Error;
		if err != nil {
			fmt.Printf("Workout Error: %+v\n",err)
			tx.Rollback()
			return tx.Error;
		}
		fmt.Printf("Workout data: %+v\n", workDay)
		fmt.Printf("Created workout id: %d\n", workDay.Workout_Id)
	}

	for g := range workDay.Groups {
		group := &workDay.Groups[g]
		group.Workout_Id = workDay.Workout_Id
		group.Workout_Day_Id = workDay.Workout_Day_Id
		repo.groupRepo.SaveGroup(group, tx)
	}
	return nil;

}