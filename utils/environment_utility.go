package utils

import (
	"github.com/joho/godotenv"
	"os"
)

type IEnvironmentUtility interface {
	GetEnvironmentString (key string) string
}

type EnvironmentUtility struct {}



func NewEnvironmentUtility() *EnvironmentUtility{
	return &EnvironmentUtility{}
}


func (util *EnvironmentUtility) GetEnvironmentString (key string) string {
	var url string
	if err:= godotenv.Load(); err != nil {
		uri,exists:=os.LookupEnv(key);
		if exists {
			url=uri
		}	
	} else{
		uri,exists:=os.LookupEnv(key);
		if exists {
			url=uri
		}
	}
	return url
}



