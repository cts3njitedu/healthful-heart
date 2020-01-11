package mongorepo

import (
	"go.mongodb.org/mongo-driver/bson"
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"context"
	"fmt"
	"log"
)



type PageRepository struct {
	connection connections.IMongoConnection
}

func NewPageRepository(conn connections.IMongoConnection) *PageRepository {
	return &PageRepository{connection: conn,}
}

func (pageRepo PageRepository) GetPage(pageType string) models.Page {
	var result models.Page
	client,err:=pageRepo.connection.GetConnection();
	if err!=nil {
		log.Println(err)
		panic(err)
	}
	fmt.Println("Making db connection...");
	db:=client.Database("HealthfulHeartConfig");
	fmt.Println("Retrieving collection...");
	collection:=db.Collection("Page")
	filter:=bson.M{"pageId": pageType};
	fmt.Printf("Retrieving %s\n",pageType);
	err=collection.FindOne(context.TODO(), filter).Decode(&result)
	if err!= nil {
		log.Println(err)
		panic(err);
	}
	return result
}

