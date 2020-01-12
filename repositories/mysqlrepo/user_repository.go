package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"errors"
	"fmt"
)

const SQL_GET_USER string = "select * from User where username = ?"
const SQL_INSERT_USER string = "insert into User(firstname,lastname,email,password,username) values (?,?,?,?,?)"



type UserRepository struct {

	connection connections.IMysqlConnection
}

func NewUserRepository(connection connections.IMysqlConnection) *UserRepository {
	return &UserRepository{connection}
}

func (userRepository *UserRepository) GetUser(user models.User) models.User  {
	var queriedUser models.User
	db, err := userRepository.connection.GetDBObject();
	
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Credentials are: %s \n", user.Username)

	defer db.Close()
	row:= db.QueryRow(SQL_GET_USER, user.Username)
	
	err=row.Scan(&queriedUser.UserId, &queriedUser.FirstName, &queriedUser.LastName, &queriedUser.Email, &queriedUser.Username, &queriedUser.Password);

	if err != nil {
		panic(err.Error())
	}

	
	return queriedUser;
}

func (userRepository *UserRepository) CreateUser(user *models.User) error {
	db, err := userRepository.connection.GetDBObject();

	if err!=nil {
		
		panic(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare(SQL_INSERT_USER);


	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(&user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Username)

	if err != nil {
		panic(err.Error())
	}
	rowCnt, err := res.RowsAffected()

	if err != nil {
		panic(err.Error())
	}

	if rowCnt != 1 {
		return errors.New("Failure in sql execution")
	}
	
	return nil

}