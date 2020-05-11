package mongorepo

import (
	"go.mongodb.org/mongo-driver/bson"
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/utils"
	"context"
	"fmt"
	"log"
	"time"
)



type PageRepository struct {
	connection connections.IMongoConnection
	environmentUtil utils.IEnvironmentUtility
}

func NewPageRepository(connection connections.IMongoConnection, environmentUtil utils.IEnvironmentUtility) *PageRepository {
	return &PageRepository{connection, environmentUtil}
}

func (pageRepo PageRepository) GetPage(pageType string) models.Page {
	var result models.Page
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client,err:=pageRepo.connection.GetConnection(ctx);
	
	if err!=nil {
		log.Println(err)
		panic(err)
	}
	fmt.Println("Making db connection...");
	dbName:=pageRepo.environmentUtil.GetEnvironmentString("MONGODB_HEALTH_CONFIG_DB")
	db:=client.Database(dbName);
	fmt.Println("Retrieving collection...");
	collection:=db.Collection("Page")
	filter:=bson.M{"pageId": pageType};
	fmt.Printf("Retrieving %s\n",pageType);
	err=collection.FindOne(ctx, filter).Decode(&result)
	if err!= nil {
		log.Println(err)
		panic(err);
	}
	err = client.Disconnect(ctx)
	if err!= nil {
		log.Println(err)
		panic(err);
	}
	return result
}

