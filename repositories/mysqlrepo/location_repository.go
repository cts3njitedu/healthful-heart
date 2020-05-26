package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	// "strings"
	Util "github.com/cts3njitedu/healthful-heart/utils"
	// "github.com/jinzhu/gorm"
)
type LocationRepository struct{
	connection connections.IMysqlConnection
}

func NewLocationRepository(connection connections.IMysqlConnection) * LocationRepository {
	return &LocationRepository{connection}
}

func (repo * LocationRepository) GetByLocationIds(ids []string) ([]models.Location, error) {
	locations := make([]models.Location, 0); 
	db, err := repo.connection.GetGormConnection();
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	db.Table("Location").Where("LOCATION_ID in (?)", ids).Find(&locations)
	return locations, nil;
}
func (repo * LocationRepository) GetLocations() ([]models.Location, error) {
	var locations []models.Location
	db, err := repo.connection.GetGormConnection();
	// defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	db.Table("Location").Find(&locations)
	return locations, nil;
}

func (repo * LocationRepository) GetLocationsQueryParams(queryOptions models.QueryOptions) ([]models.Location, error) {
	var locations []models.Location
	db, err := repo.connection.GetGormConnection();
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	columns := map[string]models.QueryOptions {
		"state" : models.QueryOptions{},
		"city" : models.QueryOptions{},
		"location_id" : models.QueryOptions{},
		"name" : models.QueryOptions{},
		"country" : models.QueryOptions{},
		"zipcode" : models.QueryOptions{},
		"location" : models.QueryOptions{},

	}

	sortMap := map[string]models.QueryOptions {
		"asc" : models.QueryOptions{},
		"desc" : models.QueryOptions{},
	}
	totalQuery, values := Util.SqlQueryBuilder(queryOptions, columns, sortMap, "Location");
	rows, err := db.Raw(totalQuery, values...).Rows()
	
	if err != nil {
		fmt.Printf("There was an error: %+v\n", err)
	} else {
		for rows.Next() {
			location := models.Location{};
			if err := db.ScanRows(rows, &location); err != nil {
				fmt.Printf("Error: %+v\n", err)
			}
			locations = append(locations, location)
		}
	}
	return locations, nil
}