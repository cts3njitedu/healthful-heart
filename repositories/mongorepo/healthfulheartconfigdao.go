package mongorepo

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"context"
	"fmt"
)


func GetPage(pageType string) models.Page {
	var result models.Page
	client,err:=connections.GetConnection();
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Making db connection...");
	db:=client.Database("HealthfulHeartConfig");
	fmt.Println("Retrieving collection...");
	collection:=db.Collection("Page")
	filter:=bson.M{"pageType": pageType};
	fmt.Println("Executing query...");
	err=collection.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Printf("The result is %+v\n",result);
	if err!= nil {
		log.Fatal(err);
	}
	return result
}

