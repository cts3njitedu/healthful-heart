package models

type HeartResponse struct {
	ActionType string `json:"actionType"`
	SectionInfos []SectionInfo `json:"sectionInfos"`
	NewSections []Section `json:"newSections"`
	Message string `json:"message"`
	IsSuccessful bool `json:"isSuccessful"`
	Data interface{} `json:"data"`
	WorkoutDays [] WorkoutDayLocationRequest `json:"workoutDays"`
	ActionInfo ActionInfo `json:"actionInfo"`
}

type ActionInfo struct {
	Added map[string][]interface{}
	Deleted map[string][]interface{}
	Modified map[string][]interface{}
}

type SectionInfo struct {
	SectionMetaData SectionMetaData `json:"sectionMetaData"`
	Section Section `json:"section"`
}

type SectionMetaData struct {
	Id string `json:"id"`
	AssociatedIds map[string]interface{} `json:"associatedIds"`
	Page int64  `json:"page"`
	VersionNb int64 `json:"versionNb"`
	TableHeaders []string `json:"tableHeaders"`
}