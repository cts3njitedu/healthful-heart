package connections

import (
	"go.mongodb.org/mongo-driver/mongo"
	"database/sql"
	"github.com/streadway/amqp"
	"github.com/jinzhu/gorm"
	"context"
)

type IMongoConnection interface {
	GetConnection(ctx context.Context) (*mongo.Client, error)
	GetFileConnection() (*mongo.Client, error)
}

type IMysqlConnection interface {
	GetDBObject() (*sql.DB, error)
	GetGormConnection() (*gorm.DB, error)
}

type IRabbitConnection interface {
	GetConnection() (*amqp.Connection, error)
}