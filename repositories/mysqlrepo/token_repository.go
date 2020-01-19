package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"errors"
	"time"
	
)

const SQL_INSERT_TOKEN string = "insert into UserToken(Refresh_token,expiration_time,user_id) values (?,?,?)"

type TokenRepository struct {
	connection connections.IMysqlConnection
}

func NewTokenRepository(connection connections.IMysqlConnection) *TokenRepository {
	return &TokenRepository{connection}
}


func (repo * TokenRepository) SaveRefreshToken(token string, expirationTime time.Time, userId string) error {
	db, err := repo.connection.GetDBObject();

	if err!=nil {
		
		panic(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare(SQL_INSERT_TOKEN);


	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(token, expirationTime, userId)

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