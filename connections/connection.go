package connections

import (
	"go.mongodb.org/mongo-driver/mongo"
	"database/sql"
)

type IMongoConnection interface {
	GetConnection() (*mongo.Client, error)
}

type IMysqlConnection interface {
	GetDBObject() (*sql.DB, error)
}