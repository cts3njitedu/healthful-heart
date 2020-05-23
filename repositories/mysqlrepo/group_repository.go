package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
	Util "github.com/cts3njitedu/healthful-heart/utils"
)

type GroupRepository struct {
	connection connections.IMysqlConnection
}

func NewGroupRepository(connection connections.IMysqlConnection) * GroupRepository {
	return &GroupRepository{connection}
}
func (repo *GroupRepository) GetGroupByParams(queryOptions models.QueryOptions) ([]models.Group, error) {
	var groups []models.Group
	db, err := repo.connection.GetGormConnection();
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	columns := map[string]models.QueryOptions {
		"workout_day_id" : models.QueryOptions{},
		"workout_id" : models.QueryOptions{},
		"group_id" : models.QueryOptions{},
		"sequence" : models.QueryOptions{},
		"sets" : models.QueryOptions{},
		"repetitions" : models.QueryOptions{},
		"weight" : models.QueryOptions{},
		"variation" : models.QueryOptions{},
		"duration" : models.QueryOptions{},
		"cre_ts" : models.QueryOptions{},
		"mod_ts" : models.QueryOptions{},
		"del_ts" : models.QueryOptions{},
		"version_nb" : models.QueryOptions{},
	}

	sortMap := map[string]models.QueryOptions {
		"asc" : models.QueryOptions{},
		"desc" : models.QueryOptions{},
	}
	totalQuery, values := Util.SqlQueryBuilder(queryOptions, columns, sortMap, "`Group`");

	rows, err := db.Raw(totalQuery, values...).Rows()

	if err != nil {
		fmt.Printf("There was an error: %+v\n", err)
	} else {
		for rows.Next() {
			group := models.Group{};
			if err := db.ScanRows(rows, &group); err != nil {
				fmt.Printf("Error: %+v\n", err)
			}
			groups = append(groups, group)
		}
	}
	return groups, nil;
}
func (repo *GroupRepository) SaveGroup(group *models.Group, tx *gorm.DB) error {
	defer func() {
		if r := recover(); r != nil {
		  tx.Rollback()
		}
	}()
	if group.Group_Id != 0 {
		ret := tx.Table("Group").
			Where("group_id = ? AND version_nb = ?",group.Group_Id, group.Version_Nb).
			Updates(map[string]interface{}{"mod_ts": time.Now(), "version_nb": group.Version_Nb + 1});
		fmt.Printf("Rows affected: %d, Group Id: %d\n",ret.RowsAffected,group.Group_Id)
	} else {
		t := time.Now()
		creTs := t.Format("2006-01-02 15:04:05")
		group.Cre_Ts = &creTs;
		group.Version_Nb = 1;
		err := tx.Table("Group").Create(&group).Error;
		if err != nil {
			fmt.Printf("Group error: %+v\n", err)
			tx.Rollback()
			return tx.Error;
		}
	
		fmt.Printf("Created group id: %d\n", group.Group_Id)
	}
	
	return nil
}