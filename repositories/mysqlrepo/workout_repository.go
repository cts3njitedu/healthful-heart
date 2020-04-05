package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
)


const SQL_QUERY_WORKOUT string = "workout_day_id = ? AND workout_type_cd = ?"
type WorkoutRepository struct {
	connection connections.IMysqlConnection
	groupRepo IGroupRepository
}

func NewWorkoutRepository(connection connections.IMysqlConnection, groupRepo IGroupRepository) * WorkoutRepository {
	return &WorkoutRepository{connection, groupRepo}
}

func (repo * WorkoutRepository) SaveWorkout(workDay *models.Workout, tx *gorm.DB) error {
	var workoutQuery *models.Workout = &models.Workout{}
	defer func() {
		if r := recover(); r != nil {
		  tx.Rollback()
		}
	}()
	tx.Table("Workout").Where(SQL_QUERY_WORKOUT, workDay.Workout_Day_Id, workDay.Workout_Type_Cd).
			First(&workoutQuery);

	if workoutQuery.Workout_Id != 0 {
		workDay.Workout_Id = workoutQuery.Workout_Id
		ret := tx.Table("Workout").
			Where("workout_id = ?",workDay.Workout_Id).
			Update("mod_ts", time.Now());
		fmt.Printf("Rows affected: %d, Workout Id: %d\n",ret.RowsAffected,workDay.Workout_Id)	
	} else {
		t := time.Now()
		creTs := t.Format("2006-01-02 15:04:05")
		workDay.Cre_Ts = &creTs;
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