package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	"strings"
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
	length := 0;
	selectQuery := []string{"*"};
	if queryOptions.Select != nil {
		selectQuery = make([]string, 0 , len(queryOptions.Select));
		for _, v := range queryOptions.Select {
			if _, ok := columns[strings.ToLower(v)]; ok {
				selectQuery = append(selectQuery, strings.ToUpper(v));
			}
		}
		
	}
	fmt.Printf("Select Query: %+v\n", selectQuery)
	whereQuery := []string{"1=?"};
	var whereValues []interface{};
	whereValues = append(whereValues, 1);
	if queryOptions.Where != nil {
		whereQuery = make([]string, 0, len(queryOptions.Where))
		whereValues = make([]interface{}, 0, len(queryOptions.Where))
		for k, v := range queryOptions.Where {
			if nv, ok := v.(string); ok {
				whereValues = append(whereValues, "%" + nv + "%")
				whereQuery = append(whereQuery, fmt.Sprintf("%s LIKE ?", k))
			} else {
				whereValues = append(whereValues, v)
				whereQuery = append(whereQuery, fmt.Sprintf("%s = ?", k))
			}
			
		}
	}
	fmt.Printf("Where Query: %+v, Where Values: %+v\n", whereQuery, whereValues)
	orderQuery := []string {"1"}
	if queryOptions.Order != nil {
		orderQuery = make([]string, 0, len(queryOptions.Order))
		for k, v := range queryOptions.Order {
			if _, ok := columns[strings.ToLower(k)]; ok {
				if _, sok := sortMap[strings.ToLower(v)]; sok {
					orderQuery = append(orderQuery, k + " " + v)
				}
			}
		}
	}
	fmt.Printf("Order Query: %+v\n", orderQuery);

	length = length + len(whereValues)
	values := make([]interface{}, 0, length)
	values = append(values, whereValues...)

	totalQuery := "SELECT " + strings.Join(selectQuery, " , ") + " FROM Location WHERE " + strings.Join(whereQuery, " AND ") +
		" ORDER BY " + fmt.Sprintf("%v", strings.Join(orderQuery, " , "));
	// totalQuery := "Select * from location"
	fmt.Println("Query:", totalQuery)
	fmt.Printf("Values: %+v\n",values)
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