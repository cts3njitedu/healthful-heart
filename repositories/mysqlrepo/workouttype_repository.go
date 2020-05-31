package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	Util "github.com/cts3njitedu/healthful-heart/utils"
)

const SQL_GET_WORKOUT_TYPES string = "select * from WorkoutType"

type WorkoutTypeRepository struct {
	connection connections.IMysqlConnection
}

func NewWorkoutTypeRepository(connection connections.IMysqlConnection) * WorkoutTypeRepository {
	return &WorkoutTypeRepository{connection}
}

func (repo * WorkoutTypeRepository) GetWorkoutTypes(queryOptions models.QueryOptions) ([]models.WorkoutType, error) {
	var workTypes []models.WorkoutType
	db, err := repo.connection.GetGormConnection();
	// defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	columns := map[string]models.QueryOptions {
		"category_cd" : models.QueryOptions{},
		"workout_type_id" : models.QueryOptions{},
		"workout_type_desc" : models.QueryOptions{},
	}

	sortMap := map[string]models.QueryOptions {
		"asc" : models.QueryOptions{},
		"desc" : models.QueryOptions{},
	}

	totalQuery, values := Util.SqlQueryBuilder(queryOptions, columns, sortMap, "WorkoutType");

	rows, err := db.Raw(totalQuery, values...).Rows()
	
	if err != nil {
		fmt.Printf("There was an error: %+v\n", err)
	} else {
		for rows.Next() {
			workType := models.WorkoutType{};
			if err := db.ScanRows(rows, &workType); err != nil {
				fmt.Printf("Error: %+v\n", err)
			}
			workTypes = append(workTypes, workType)
		}
	}
	return workTypes, nil
}

