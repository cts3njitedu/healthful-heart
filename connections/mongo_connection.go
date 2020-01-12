package connections

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"fmt"
	"log"
	"github.com/cts3njitedu/healthful-heart/utils"
)

type MongoConnection struct {
	environmentUtil utils.IEnvironmentUtility

}

func NewMongoConnection(environmentUtil utils.IEnvironmentUtility) *MongoConnection {
	return &MongoConnection{environmentUtil}
}

func (m *MongoConnection) GetConnection() (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	url:=m.environmentUtil.GetEnvironmentString("MONGODB_URI")
	clientOptions:=options.Client().ApplyURI(url)
	
	client,err := mongo.Connect(ctx, clientOptions)
	if err!= nil {
		log.Println(err)
		panic(err)
	}
	err = client.Ping(ctx,nil)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client,err
}