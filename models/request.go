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
	LocationId string
	WorkoutDayId string
	WorkoutId string
	SectionInfo `json:"sectionInfo"`
	WorkoutDays [] WorkoutDayLocationRequest

}

type WorkoutDayLocationRequest struct {
	WorkoutDayId string `json:"workoutDayId"`
	Version_Nb string `json:"versionNb"`
	IsDeleted bool `json:"isDeleted"`
	Fields [] Field `json:"fields"`
	Workouts [] WorkoutRequest `json:"workouts"`
}

type WorkoutRequest struct {
	WorkoutId string `json:"workoutId"`
	Version_Nb string `json:"versionNb"`
	IsDeleted bool `json:"isDeleted"`
	Fields []Field `json:"fields"`
	Groups []GroupRequest `json:"groups"`

}

type GroupRequest struct {
	GroupId string `json:"groupId"`
	Version_Nb string `json:"versionNb"`
	IsDeleted bool `json:"isDeleted"`
	Fields []Field `json:"fields"`
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
	VIEW_WORKOUTS_HEADER = "VIEW_WORKOUTS_HEADER"
	VIEW_WORKOUTS = "VIEW_WORKOUTS"
	VIEW_WORKOUT_DETAILS_META_INFO = "VIEW_WORKOUT_DETAILS_META_INFO"
	VIEW_WORKOUT_DETAILS = "VIEW_WORKOUT_DETAILS"
	WORKOUTS_ACTION_ERRORS = "WORKOUTS_ACTION_ERRORS"
	WORKOUTS_ACTION_SUCCESS = "WORKOUTS_ACTION_SUCCESS"
	WORKOUTS_ACTION = "WORKOUTS_ACTION"
	ACTION_WORKOUT = "ACTION_WORKOUT"
	ACTION_WORKOUTDAY_LOCATION = "ACTION_WORKOUTDAY_LOCATION"
)
