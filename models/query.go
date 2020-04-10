package models


type QueryOptions struct {
	Where map[string]interface{}
	Order string
	In map[string][]interface{}
}