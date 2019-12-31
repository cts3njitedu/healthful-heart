package models


type Page struct {
	PageType string `bson:"pageType" json:"pageType"`
	Sections  []Section `bson:"sections" json:"sections"`
	Errors []string `bson:"errors" json:"errors"`
}




type Section struct {
	SectionType string `bson:"sectionType" json:"sectionType"`
	Fields []Field	`bson:"fields" json:"fields"`
	Errors []string `bson:"errors" json:"errors"`
}


type Field struct {
	FieldType string `bson:"fieldType" json:"fieldType"`
	Name string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
	Placeholder string `bson:"placeHolder" json:"placeHolder"`
	IsDisabled bool `bson:"isDisabled" json:"isDisabled"`
	Type string `bson:"type" json:"type"`
	IsHidden bool `bson:"isHidden" json:"isHidden"`
	Items []Item `bson:"items" json:"items"`
	IsDirty bool `bson:"isDirty" json:"isDirty"`
	Errors []string `bson:"errors" json:"errors"`
	Validations []Validation `bson:"validations" json:"validations"`
}

type Item struct {
	Item string `bson:"item" json:"item"`
}

type Validation struct {
	ValidationName string `bson:"validationName" json:"validationName"`
	IsMandatory bool `bson:"isMandatory" json:"isMandatory"`
	RegexValue string `bson:"regexValue" json:"regexValue"`
	Message string `bson:"message" json:"message"`
	MinLength int `bson:"minLength" json:"minLength"`
	MaxLength int `bson:"maxLength" json:"maxLength"`
	IsValid bool 
}