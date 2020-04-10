package models

type HeartResponse struct {
	SectionInfos []SectionInfo `json:"sectionInfos"`
	NewSections []Section `json:"newSections"`
}

type SectionInfo struct {
	SectionMetaData SectionMetaData `json:"sectionMetaData"`
	Section Section `json:"section"`
}

type SectionMetaData struct {
	Id string `json:"id"`
}