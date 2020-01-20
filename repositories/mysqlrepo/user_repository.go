package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"errors"
	"fmt"
	"database/sql"
	"strconv"
)

const SQL_GET_USER string = "select * from User where username = ?"
const SQL_INSERT_USER string = "insert into User(firstname,lastname,email,password,username) values (?,?,?,?,?)"
const SQL_GET_TOKEN string = "select u.User_id, u.Firstname, u.LastName, u.Email, u.Username, ut.Refresh_token, ut.Expiration_Time" +
		" from User u inner join UserToken ut on u.user_id = ut.user_id where u.user_id = ?"


type UserRepository struct {

	connection connections.IMysqlConnection
}

func NewUserRepository(connection connections.IMysqlConnection) *UserRepository {
	return &UserRepository{connection}
}

type UserError struct {
	s string
}

func (userError * UserError) Error() string {
	return userError.s
}

func (repo * UserRepository) GetUserToken(userId string) (models.User, error) {
	var queriedUser models.User
	db, err := repo.connection.GetDBObject();
	
	if err != nil {
		panic(err.Error())
	}


	defer db.Close()
	row:= db.QueryRow(SQL_GET_TOKEN, userId)
	
	err=row.Scan(&queriedUser.User_Id,&queriedUser.FirstName, &queriedUser.LastName, &queriedUser.Email, &queriedUser.Username, &queriedUser.RefreshToken, &queriedUser.ExpirationTime);

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("User Token is missing")
		}
		panic(err.Error())
	}
	
	return queriedUser, nil;
}

func (userRepository *UserRepository) GetUser(user models.User) (models.User, error) {
	var queriedUser models.User
	db, err := userRepository.connection.GetDBObject();
	
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Credentials are: %s \n", user.Username)

	defer db.Close()
	row:= db.QueryRow(SQL_GET_USER, user.Username)
	
	err=row.Scan(&queriedUser.User_Id, &queriedUser.FirstName, &queriedUser.LastName, &queriedUser.Email, &queriedUser.Username, &queriedUser.Password);

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, &UserError{"Invalid username or password"}
		}
		panic(err.Error())
	}
	
	return queriedUser, nil;
}

func (userRepository *UserRepository) CreateUser(user *models.User) (error) {
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

	id, err := res.LastInsertId();
	
	if err != nil {
		panic(err.Error())
	}
	user.User_Id = strconv.FormatInt(id, 10)
	return nil

}