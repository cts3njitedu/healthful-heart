package connections

import (
	"fmt"
    "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/cts3njitedu/healthful-heart/utils"
)

type MysqlConnection struct {
	environmentUtil utils.IEnvironmentUtility
}



func NewMysqlConnection(environmentUtil utils.IEnvironmentUtility) *MysqlConnection {
	return &MysqlConnection{environmentUtil}
}

func (conn *MysqlConnection) GetDBObject() (*sql.DB, error) {
	fmt.Println("Creating mysql db object");
	var url string;
	if err:= godotenv.Load(); err != nil {
		uri,exists:=os.LookupEnv("CLEARDB_DATABASE_URL");
		if exists {
			url=uri
		}	
	} else{
		uri,exists:=os.LookupEnv("CLEARDB_DATABASE_URL");
		if exists {
			url=uri
		}
	}

	db, err := sql.Open("mysql", url);

	if err != nil {
		fmt.Println("This is the error",err);
		panic(err.Error())
	}
	return db, err
}

func (conn *MysqlConnection) GetGormConnection() (*gorm.DB, error) {
	fmt.Println("Getting gorm connection...")
	url := conn.environmentUtil.GetEnvironmentString("CLEARDB_DATABASE_URL")
	db, err := gorm.Open("mysql", url)
	if err != nil {
		fmt.Println("This is the error",err);
		panic(err.Error())
	}
	return db, err
}