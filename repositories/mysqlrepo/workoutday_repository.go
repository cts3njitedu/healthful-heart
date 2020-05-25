package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
	"strings"
	Util "github.com/cts3njitedu/healthful-heart/utils"
	"errors"
)

const SQL_QUERY_WORKOUT_DAY string = "user_id = ? AND location_id = ? and workout_date = ?"
const SQL_QUERY_WORKOUT_DAYS_USER_ID string = "user_id = ?"
type WorkoutDayRepository struct {
	connection connections.IMysqlConnection
	workoutRepo IWorkoutRepository
	groupRepo IGroupRepository
}

func NewWorkoutDayRepository(connection connections.IMysqlConnection, workoutRepo IWorkoutRepository, groupRepo IGroupRepository) * WorkoutDayRepository {
	return &WorkoutDayRepository{connection, workoutRepo, groupRepo}
}

func (repo * WorkoutDayRepository) DeleteWorkoutDays(ids map[string][]string, tx *gorm.DB) bool {
	ret := tx.Table("WorkoutDay").Where("Workout_Day_Id IN (?)", ids["WorkoutDay"]).Delete(models.WorkoutDay{})
	if ret.Error != nil {
		fmt.Printf("Unable to delete Workouts %+v\n", ret.Error)
		return false;
	} else {
		fmt.Printf("Workout Days Deleted: %+v\n", ret.RowsAffected)
	}
	return true;
}
func (repo *WorkoutDayRepository) GetWorkoutDaysLocationByParams(queryOptions models.QueryOptions) ([]models.WorkoutDay, error) {
	var workoutDays []models.WorkoutDay
	db, err := repo.connection.GetGormConnection();
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	columns := map[string]models.QueryOptions {
		"workout_day_id" : models.QueryOptions{},
		"user_id" : models.QueryOptions{},
		"location_id" : models.QueryOptions{},
		"workout_date" : models.QueryOptions{},
		"cre_ts" : models.QueryOptions{},
		"mod_ts" : models.QueryOptions{},
		"del_ts" : models.QueryOptions{},
		"version_nb" : models.QueryOptions{},
		"workout_file_id" : models.QueryOptions{},
		"name" : models.QueryOptions{},
		"state" : models.QueryOptions{},
		"city" : models.QueryOptions{},
		"country" : models.QueryOptions{},
		"zipcode" : models.QueryOptions{},
		"location" : models.QueryOptions{},

	}

	sortMap := map[string]models.QueryOptions {
		"asc" : models.QueryOptions{},
		"desc" : models.QueryOptions{},
	}
	totalQuery, values := Util.SqlQueryBuilder(queryOptions, columns, sortMap, "View_Workout_Day_Location");

	rows, err := db.Raw(totalQuery, values...).Rows()
	
	if err != nil {
		fmt.Printf("There was an error: %+v\n", err)
	} else {
		
		for rows.Next() {
			workoutDay := models.WorkoutDay{};
			location := models.Location{};
			err = rows.Scan(&workoutDay.Workout_Day_Id, &workoutDay.Location_Id, &workoutDay.Version_Nb, &location.Name, 
				&location.State, &location.City, &location.Country, &location.Zipcode, &location.Location)
			workoutDay.Location = location
			workoutDays = append(workoutDays, workoutDay)
		}
	}
	return workoutDays, nil;
}
func (repo *WorkoutDayRepository) GetWorkoutDaysSpecifyColumns(queryOptions models.QueryOptions) ([]models.WorkoutDay, error) {
	var workoutDays []models.WorkoutDay
	db, err := repo.connection.GetGormConnection();
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	var whereClause map[string]interface{};
	if queryOptions.Where != nil {
		whereClause = queryOptions.Where;
	}
	selectClause := "*"
	if queryOptions.Select != nil {
		selectClause = strings.Join(queryOptions.Select, ",")
		selectClause = strings.ToUpper(selectClause)
	}
	fmt.Printf("query: %+v\n", selectClause)
	db.Table("WorkoutDay").Select(selectClause).Where(whereClause).Order("WORKOUT_DATE").Find(&workoutDays)
	return workoutDays, nil;

} 
func (repo *WorkoutDayRepository) GetWorkoutDaysByParams(queryOptions models.QueryOptions) ([]models.WorkoutDay, error) {
	var workoutDays []models.WorkoutDay
	db, err := repo.connection.GetGormConnection();
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	columns := map[string]models.QueryOptions {
		"workout_day_id" : models.QueryOptions{},
		"user_id" : models.QueryOptions{},
		"location_id" : models.QueryOptions{},
		"workout_date" : models.QueryOptions{},
		"cre_ts" : models.QueryOptions{},
		"mod_ts" : models.QueryOptions{},
		"del_ts" : models.QueryOptions{},
		"version_nb" : models.QueryOptions{},
		"workout_file_id" : models.QueryOptions{},

	}

	sortMap := map[string]models.QueryOptions {
		"asc" : models.QueryOptions{},
		"desc" : models.QueryOptions{},
	}
	totalQuery, values := Util.SqlQueryBuilder(queryOptions, columns, sortMap, "WorkoutDay");

	rows, err := db.Raw(totalQuery, values...).Rows()
	
	if err != nil {
		fmt.Printf("There was an error: %+v\n", err)
	} else {
		for rows.Next() {
			workoutDay := models.WorkoutDay{};
			if err := db.ScanRows(rows, &workoutDay); err != nil {
				fmt.Printf("Error: %+v\n", err)
			}
			workoutDays = append(workoutDays, workoutDay)
		}
	}
	return workoutDays, nil;

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

func (repo * WorkoutDayRepository) SaveWorkoutDayLocation(workDay *models.WorkoutDay) (*models.WorkoutDay, error) {
	db, err := repo.connection.GetGormConnection();
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	

	t := time.Now()
	creTs := t.Format("2006-01-02 15:04:05")
	workDay.Cre_Ts = &creTs;
	workDay.Version_Nb = 1;
	err = db.Table("WorkoutDay").Where(SQL_QUERY_WORKOUT_DAY, workDay.User_Id, 
		workDay.Location_Id, workDay.Workout_Date).FirstOrCreate(&workDay).Error;
	if err != nil {
		fmt.Printf("WorkoutDay Error: %+v\n",err)
		return &models.WorkoutDay{}, err;
	}
	fmt.Printf("Created workout day id: %d\n", workDay.Workout_Day_Id)
	return workDay, nil
}

func (repo *WorkoutDayRepository) UpdateAllWorkoutDay(workDays []models.WorkoutDay, ids map[string][]string) error {
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
		
		for w := range workDays {
			workDay := &workDays[w]
			err = repo.SaveWorkoutDay(workDay, tx)
			if err != nil {
				tx.Rollback()
				fmt.Printf("Workout Day Save Error: %+v\n", workDay.Workout_Day_Id)
				return err
			}
		}

		if ids == nil {
			return nil
		}
		ret := repo.groupRepo.DeleteGroups(ids, tx)
		if !ret {
			tx.Rollback();
			return errors.New("Unable to Delete Group\n")
		}

		ret = repo.workoutRepo.DeleteWorkouts(ids, tx)
		if !ret {
			tx.Rollback();
			return errors.New("Unable to Delete Workouts\n")
		}

		ret = repo.DeleteWorkoutDays(ids, tx);
		if !ret {
			tx.Rollback();
			return errors.New("Unable to Delete Workout Days\n")
		}


		return nil
	})
}

func (repo * WorkoutDayRepository) SaveWorkoutDay(workDay *models.WorkoutDay, tx *gorm.DB) error {
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var workDayQuery *models.WorkoutDay = &models.WorkoutDay{}

	tx.Table("WorkoutDay").Where(SQL_QUERY_WORKOUT_DAY, workDay.User_Id, workDay.Location_Id, workDay.Workout_Date).
	First(&workDayQuery)
	if workDay.Workout_Day_Id != 0 {
		if ((workDayQuery.Workout_Day_Id !=0 && 
			(workDayQuery.Workout_Day_Id != workDay.Workout_Day_Id))) {
				tx.Rollback();
				fmt.Printf("Incorrect info for Workout Day: %+v\n", workDay.Workout_Day_Id)
				return errors.New(fmt.Sprintf("Incorrect info for Workout Day: %+v\n", workDay.Workout_Day_Id))
		}
		workDayQuery = workDay
		workDayQuery.Version_Nb = workDay.Version_Nb
	}

	if workDayQuery.Workout_Day_Id != 0 {
		workDay.Workout_Day_Id = workDayQuery.Workout_Day_Id;
		ret := tx.Table("WorkoutDay").
			Where("workout_day_id = ? AND version_nb = ?",workDay.Workout_Day_Id, workDayQuery.Version_Nb).
			Updates(map[string]interface{}{"mod_ts": time.Now(), "version_nb": workDayQuery.Version_Nb + 1});
		fmt.Printf("Rows affected: %d, Workout Day Id: %d\n",ret.RowsAffected,workDay.Workout_Day_Id)
		if (ret.RowsAffected == 0) {
			tx.Rollback()
			errString := fmt.Sprintf("Unable to Find Workout Day: %+v\n", workDay.Workout_Day_Id)
			fmt.Print(errString)
			return errors.New(errString)
		}
	
	} else {
		t := time.Now()
		creTs := t.Format("2006-01-02 15:04:05")
		workDay.Cre_Ts = &creTs;
		workDay.Mod_Ts = nil;
		workDay.Version_Nb = 1;
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
			tx.Rollback();
			return err;
		}
	}

	return nil;

}
