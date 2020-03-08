package services

import (
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/models"
	"errors"
	"fmt"
)

type LocationService struct {
	locationRepository mysqlrepo.ILocationRepository
}
var locationMap = make(map[int64]models.Location);
func NewLocationService(locationRepository mysqlrepo.ILocationRepository) *LocationService {
	locations, _ := locationRepository.GetLocations()
	fmt.Printf("Locations: %+v", locations)
	for _, location := range locations {
		locationMap[location.Location_Id] = location;
	}
	return &LocationService{locationRepository}
}

func(serv * LocationService) GetLocation(locationId int64) (models.Location, error) {
	if (locationMap[locationId] != models.Location{}) {
		return locationMap[locationId], nil
	}
	return models.Location{}, errors.New("Unable to find location")
}