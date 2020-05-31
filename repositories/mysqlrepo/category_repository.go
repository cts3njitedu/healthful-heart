package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	Util "github.com/cts3njitedu/healthful-heart/utils"
)



type CategoryRepository struct {
	connection connections.IMysqlConnection
}

func NewCategoryRepository(connection connections.IMysqlConnection) * CategoryRepository {
	return &CategoryRepository{connection}
}

const SQL_CATEGORY_WORKOUT_TYPE string = "(SELECT C.*, W.WORKOUT_TYPE_ID, W.WORKOUT_TYPE_DESC FROM Category C LEFT JOIN WorkoutType W ON C.CATEGORY_CD=W.CATEGORY_CD) AS T"

func (repo * CategoryRepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	db, err := repo.connection.GetGormConnection();
	fmt.Printf("Type of db connection: %T", db)
	// defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	db.Table("Category").Order("CATEGORY_NAME").Find(&categories) 
	return categories, nil;
}
func (repo * CategoryRepository) GetCategoriesByParams(queryOptions models.QueryOptions) ([]models.Category, error) {
	var categories []models.Category
	db, err := repo.connection.GetGormConnection();
	// defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	columns := map[string]models.QueryOptions {
		"category_cd" : models.QueryOptions{},
		"category_name" : models.QueryOptions{},
	}

	sortMap := map[string]models.QueryOptions {
		"asc" : models.QueryOptions{},
		"desc" : models.QueryOptions{},
	}

	totalQuery, values := Util.SqlQueryBuilder(queryOptions, columns, sortMap, "Category");

	rows, err := db.Raw(totalQuery, values...).Rows()
	
	if err != nil {
		fmt.Printf("There was an error: %+v\n", err)
	} else {
		for rows.Next() {
			category := models.Category{};
			if err := db.ScanRows(rows, &category); err != nil {
				fmt.Printf("Error: %+v\n", err)
			}
			categories = append(categories, category)
		}
	}
	return categories, nil

}
func (repo * CategoryRepository) GetCategoriesAndWorkoutTypes(queryOptions models.QueryOptions) ([]models.Category, error) {
	var categories []models.Category
	db, err := repo.connection.GetGormConnection();
	// defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	columns := map[string]models.QueryOptions {
		"category_cd" : models.QueryOptions{},
		"category_name" : models.QueryOptions{},
		"workout_type_id" : models.QueryOptions{},
		"workout_type_desc" : models.QueryOptions{},
	}

	sortMap := map[string]models.QueryOptions {
		"asc" : models.QueryOptions{},
		"desc" : models.QueryOptions{},
	}
	totalQuery, values := Util.SqlQueryBuilder(queryOptions, columns, sortMap, SQL_CATEGORY_WORKOUT_TYPE);

	rows, err := db.Raw(totalQuery, values...).Rows()
	
	if err != nil {
		fmt.Printf("There was an error: %+v\n", err)
	} else {
		for rows.Next() {
			category := models.Category{};
			if err := db.ScanRows(rows, &category); err != nil {
				fmt.Printf("Error: %+v\n", err)
			}
			categories = append(categories, category)
		}
	}
	return categories, nil
}