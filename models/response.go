package models

type HeartResponse struct {
	ActionType string `json:"actionType"`
	SectionInfos []SectionInfo `json:"sectionInfos"`
	NewSections []Section `json:"newSections"`
	Message string `json:"message"`
	IsSuccessful bool `json:"isSuccessful"`
	Data interface{} `json:"data"`
}

type SectionInfo struct {
	SectionMetaData SectionMetaData `json:"sectionMetaData"`
	Section Section `json:"section"`
}

type SectionMetaData struct {
	Id string `json:"id"`
	Page int64  `json:"page"`
	TableHeaders []string `json:"tableHeaders"`
}