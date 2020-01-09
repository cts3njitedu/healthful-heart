package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
)

const SQL_GET_USER string = "select * from User where username = ?"

func GetUser(username string) models.User {
	var user models.User
	db, err := connections.GetDBObject();
	
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	err = db.QueryRow(SQL_GET_USER, username).Scan(&user.USER_ID, &user.FIRSTNAME, &user.LASTNAME, &user.EMAIL, &user.USERNAME, &user.PASSWORD);

	if err != nil {
		panic(err.Error())
	}

	return user;
}