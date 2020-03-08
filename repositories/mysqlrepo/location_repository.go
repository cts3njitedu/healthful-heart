package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
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