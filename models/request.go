package models

type QueryParams struct {
	Date string
	LocationIds []int64
	Filter map[string]interface{}
}

type HeartRequest struct {
	ActionType string
	HeartSort map[string]HeartSort
	HeartFilter map[string]interface{}
	Date string

}

type HeartSort struct {
	SortOrder string
}

type HeartFilter struct {
	Column string
	Value string
}

const (
	VIEW_LOCATIONS = "VIEW_LOCATIONS"
	VIEW_WORKDATE_LOCATIONS = "VIEW_WORKDATE_LOCATIONS"
	VIEW_NON_WORKDATE_LOCATIONS = "VIEW_NON_WORKDATE_LOCATIONS"
	ADD_LOCATION = "ADD_LOCATION"
)