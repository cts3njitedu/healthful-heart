package models


type QueryOptions struct {
	Where map[string]interface{}
	Order map[string]string
	In map[string]string
	Select []string
	NotIn map[string]string
	Limit int64
	Offset int64
	IsEqual bool
}