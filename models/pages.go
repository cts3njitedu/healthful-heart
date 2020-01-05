package models


type Page struct {
	PageId string `bson:"pageId" json:"pageId"`
	PageType string `bson:"pageType" json:"pageType"`
	Sections  []Section `bson:"sections" json:"sections"`
	Errors []string `bson:"errors" json:"errors"`
}




type Section struct {
	SectionId string `bson:"sectionId" json:"sectionId"`
	Fields []Field	`bson:"fields" json:"fields"`
	Errors []string `bson:"errors" json:"errors"`
}


type Field struct {
	FieldId string `bson:"fieldId" json:"fieldId"`
	Name string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
	Placeholder string `bson:"placeholder" json:"placeholder"`
	IsDisabled bool `bson:"isDisabled" json:"isDisabled"`
	Type string `bson:"type" json:"type"`
	IsHidden bool `bson:"isHidden" json:"isHidden"`
	Items []Item `bson:"items" json:"items"`
	IsDirty bool `bson:"isDirty" json:"isDirty"`
	Errors []string `bson:"errors" json:"errors"`
	Validations []Validation `bson:"validations" json:"validations"`
	MinLength int `bson:"minLength" json:"minLength"`
	MaxLength int `bson:"maxLength" json:"maxLength"`
	IsMandatory bool `bson:"isMandatory" json:"isMandatory"`
	RegexValue string `bson:"regexValue" json:"regexValue"`
	Title string `bson:"title" json:"title"`
}

type Item struct {
	ItemId string `bson:"itemId" json:"itemId"`
	Item string `bson:"item" json:"item"`
}

type Validation struct {
	ValidationId string `bson:"validationId" json:"validationId"`
	ValidationName string `bson:"validationName" json:"validationName"`
	Message string `bson:"message" json:"message"`
	IsValid bool 
}