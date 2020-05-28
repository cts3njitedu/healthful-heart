package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
)

const SQL_GET_WORKOUT_TYPES string = "select * from WorkoutType"

type WorkoutTypeRepository struct {
	connection connections.IMysqlConnection
}

func NewWorkoutTypeRepository(connection connections.IMysqlConnection) * WorkoutTypeRepository {
	return &WorkoutTypeRepository{connection}
}

func (repo * WorkoutTypeRepository) GetWorkoutTypes() ([]models.WorkoutType, error) {
	var workTypes []models.WorkoutType
	db, err := repo.connection.GetGormConnection();
	// defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	db.Table("WorkoutType").Order("name").Find(&workTypes)
	return workTypes, nil;
}