package connections

import (
	"go.mongodb.org/mongo-driver/mongo"
	"database/sql"
	"github.com/streadway/amqp"
	"github.com/jinzhu/gorm"
)

type IMongoConnection interface {
	GetConnection() (*mongo.Client, error)
	GetFileConnection() (*mongo.Client, error)
}

type IMysqlConnection interface {
	GetDBObject() (*sql.DB, error)
	GetGormConnection() (*gorm.DB, error)
}

type IRabbitConnection interface {
	GetConnection() (*amqp.Connection, error)
}