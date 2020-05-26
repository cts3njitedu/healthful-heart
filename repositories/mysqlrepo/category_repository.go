package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
)



type CategoryRepository struct {
	connection connections.IMysqlConnection
}

func NewCategoryRepository(connection connections.IMysqlConnection) * CategoryRepository {
	return &CategoryRepository{connection}
}

func (repo * CategoryRepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	db, err := repo.connection.GetGormConnection();
	fmt.Printf("Type of db connection: %T", db)
	// defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	db.Table("Category").Find(&categories)
	return categories, nil;
}