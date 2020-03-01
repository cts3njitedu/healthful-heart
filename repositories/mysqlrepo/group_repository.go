package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
)

type GroupRepository struct {
	connection connections.IMysqlConnection
}

func NewGroupRepository(connection connections.IMysqlConnection) * GroupRepository {
	return &GroupRepository{connection}
}

func (repo *GroupRepository) SaveGroup(group *models.Group, tx *gorm.DB) error {
	defer func() {
		if r := recover(); r != nil {
		  tx.Rollback()
		}
	}()
	t := time.Now()
	creTs := t.Format("2006-01-02 15:04:05")
	group.Cre_Ts = &creTs;
	tx.Table("Group").Create(&group);
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error;
	}

	fmt.Printf("Created group id: %d\n", group.Group_Id)
	return nil
}