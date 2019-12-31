package mongo_client

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/bson"
)

type Name struct {
	FirstName string `bson:"firstName" json:"firstName"`
	LastName string	`bson:"lastName" json:"lastName"`
}

var result Name

func GetConnection() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// clientOptions:=options.Client().ApplyURI("mongodb://healthconfig:12345@localhost:27017")
	clientOptions:=options.Client().ApplyURI("mongodb://localhost:27017/test")
	client,err := mongo.Connect(ctx, clientOptions)
	if err!= nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx,nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	db:=client.Database("test");
	collection:=db.Collection("mycollection")
	filter:=bson.M{"firstName":"John"}
	err=collection.FindOne(context.TODO(), filter).Decode(&result)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println(result);
	fmt.Printf("Found a single document: %+v\n",result);
	defer cancel()
}