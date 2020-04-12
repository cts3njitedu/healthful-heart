package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"strings"
	"fmt"
)
type LocationRepository struct{
	connection connections.IMysqlConnection
}

func NewLocationRepository(connection connections.IMysqlConnection) * LocationRepository {
	return &LocationRepository{connection}
}

func (repo * LocationRepository) GetLocations() ([]models.Location, error) {
	var locations []models.Location
	db, err := repo.connection.GetGormConnection();
	defer db.Close()
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
	var whereClause map[string]interface{};
	if queryOptions.Where != nil {
		whereClause = queryOptions.Where;
	}
	orderClause := "LOCATION_ID"

	if queryOptions.Order != nil {
		orderClause = strings.Join(queryOptions.Order, ",")
		orderClause = strings.ToUpper(orderClause)
	}

	fmt.Printf("Order clause: %+v\n", orderClause)

	db.Table("Location").Where(whereClause).Order(orderClause).Find(&locations)
	return locations, nil
}