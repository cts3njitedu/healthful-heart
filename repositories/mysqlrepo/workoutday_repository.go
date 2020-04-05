package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
)

const SQL_QUERY_WORKOUT_DAY string = "user_id = ? AND location_id = ? and workout_date = ?"
const SQL_QUERY_WORKOUT_DAYS_USER_ID string = "user_id = ?"
type WorkoutDayRepository struct {
	connection connections.IMysqlConnection
	workoutRepo IWorkoutRepository
}

func NewWorkoutDayRepository(connection connections.IMysqlConnection, workoutRepo IWorkoutRepository) * WorkoutDayRepository {
	return &WorkoutDayRepository{connection, workoutRepo}
}

func (repo * WorkoutDayRepository) GetWorkoutDays(userId string) ([]models.WorkoutDay, error) {
	var workoutDays []models.WorkoutDay
	db, err := repo.connection.GetGormConnection();
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	db.Table("WorkoutDay").Where(SQL_QUERY_WORKOUT_DAYS_USER_ID, userId).Order("WORKOUT_DATE").Find(&workoutDays)
	return workoutDays, nil;
}

func (repo * WorkoutDayRepository) SaveWorkoutDay(workDay *models.WorkoutDay) error {
	db, err := repo.connection.GetGormConnection();
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	return db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
			  tx.Rollback()
			}
		}()
		var workDayQuery *models.WorkoutDay = &models.WorkoutDay{}
		tx.Table("WorkoutDay").Where(SQL_QUERY_WORKOUT_DAY, workDay.User_Id, workDay.Location_Id, workDay.Workout_Date).
			First(&workDayQuery)

		if workDayQuery.Workout_Day_Id != 0 {
			workDay.Workout_Day_Id = workDayQuery.Workout_Day_Id;
			ret := tx.Table("WorkoutDay").
				Where("workout_day_id = ?",workDay.Workout_Day_Id).
				Update("mod_ts", time.Now());
			fmt.Printf("Rows affected: %d, Workout Day Id: %d\n",ret.RowsAffected,workDay.Workout_Day_Id)
		
		} else {
			t := time.Now()
			creTs := t.Format("2006-01-02 15:04:05")
			workDay.Cre_Ts = &creTs;
			err := tx.Table("WorkoutDay").Create(&workDay).Error;
			if err != nil {
				fmt.Printf("WorkoutDay Error: %+v\n",err)
				tx.Rollback()
				return tx.Error;
			}
			fmt.Printf("Created workout day id: %d\n", workDay.Workout_Day_Id)
		}

		for w := range workDay.Workouts {
			workOut := &workDay.Workouts[w]
			workOut.Workout_Day_Id = workDay.Workout_Day_Id
			err := repo.workoutRepo.SaveWorkout(workOut, tx)
			if err != nil {
				fmt.Printf("The workout error: %+v\n", workOut)
			}
		}

		return nil
	})
}
