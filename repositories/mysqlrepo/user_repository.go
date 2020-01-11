package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"errors"
)

const SQL_GET_USER string = "select * from User where username = ?"
const SQL_INSERT_USER string = "insert intto user(firstname,lastname,email,password,username) values (?,?,?,?,?)"



type UserRepository struct {

	connection connections.IMysqlConnection
}

func NewUserRepository(connection connections.IMysqlConnection) *UserRepository {
	return &UserRepository{connection}
}

func (userRepository *UserRepository) GetUser(username string) models.User {
	var user models.User
	db, err := userRepository.connection.GetDBObject();
	
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	err = db.QueryRow(SQL_GET_USER, username).Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password);

	if err != nil {
		panic(err.Error())
	}

	return user;
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

	rowCnt, err := res.RowsAffected()

	if err != nil {
		panic(err.Error())
	}

	if rowCnt != 1 {
		return errors.New("Failure in sql execution")
	}
	
	return nil

}