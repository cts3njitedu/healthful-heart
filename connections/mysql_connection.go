package connections

import (
	"fmt"
    "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
)

func GetDBObject() (*sql.DB, error) {
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
		panic(err.Error())
	}
	return db, err
}