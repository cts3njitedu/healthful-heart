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
	HeartPagination HeartPagination
	Date string
	SectionInfo `json:"sectionInfo"`

}

type HeartSort struct {
	SortOrder string
}

type HeartFilter struct {
	Column string
	Value string
}

type HeartPagination struct {
	Limit int64
	Offset int64
	Page int64
}

const (
	VIEW_LOCATIONS = "VIEW_LOCATIONS"
	VIEW_WORKOUTDATE_LOCATIONS = "VIEW_WORKOUTDATE_LOCATIONS"
	VIEW_NON_WORKOUTDATE_LOCATIONS = "VIEW_NON_WORKOUTDATE_LOCATIONS"
	ADD_WORKOUTDATE_LOCATION = "ADD_WORKOUTDATE_LOCATION"
	DEFAULT_ACTIVITY = "DEFAULT_ACTIVITY"
	LOCATION_SELECTED = "LOCATION_SELECTED"
)