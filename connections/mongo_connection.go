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
var ctx context.Context;
var cancel context.CancelFunc
var clientOptions *options.ClientOptions
func NewMongoConnection(environmentUtil utils.IEnvironmentUtility) *MongoConnection {
	return &MongoConnection{environmentUtil}
}


func (m *MongoConnection) GetConnection() (*mongo.Client, error) {
	return m.makeConnection("MONGODB_URI")
}

func (m *MongoConnection) makeConnection(connUrl string) (*mongo.Client,error) {
	fmt.Println("Creating Mongo Context...")
	// d := time.Now().Add(10*time.Second)
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	url:= m.environmentUtil.GetEnvironmentString(connUrl);
	clientOptions =options.Client().ApplyURI(url)
	helper := uint64(1)
	clientOptions.MaxPoolSize = &helper
	fmt.Printf("Client Options:%+v\n", clientOptions)
	client,err := mongo.Connect(ctx, clientOptions)
	if err!= nil {
		log.Println("sugar Honey Ice err");
		log.Println(err)
		panic(err)
	}
	err = client.Ping(ctx,nil)
	if err != nil {
		log.Println("sugar Honey Ice err1");
		log.Println(err)
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client,err
}
func (m *MongoConnection) GetFileConnection() (*mongo.Client, error) {

	return m.makeConnection("MONGODB_FILE_URI");
}