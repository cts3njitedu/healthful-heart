package connections

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"fmt"
	"os"
	"log"
	"github.com/joho/godotenv"
)

type MongoConnection struct {


}

func NewMongoConnection() *MongoConnection {
	return &MongoConnection{}
}

func (m *MongoConnection) GetConnection() (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var url string;
	if err:= godotenv.Load(); err != nil {
		uri,exists:=os.LookupEnv("MONGODB_URI");
		if exists {
			url=uri
		}	
	} else{
		uri,exists:=os.LookupEnv("MONGODB_URI");
		if exists {
			url=uri
		}
	}
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