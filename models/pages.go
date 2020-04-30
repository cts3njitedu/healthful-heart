package models


type Page struct {
	PageId string `bson:"pageId" json:"pageId"`
	PageType string `bson:"pageType" json:"pageType"`
	Sections  []Section `bson:"sections" json:"sections"`
	Errors []string `bson:"errors" json:"errors"`
}




type Section struct {
	Id string `bson:"id" json:"id"`
	ParentId string `bson:"parentId" json:"parentId"`
	SectionId string `bson:"sectionId" json:"sectionId"`
	Fields []Field	`bson:"fields" json:"fields"`
	Errors []string `bson:"errors" json:"errors"`
	IsHidden bool `bson:"isHidden" json:"isHidden"`
}


func (sect * Section) FindField(fieldId string, section Section) Field {
	for f := range section.Fields {
		field := section.Fields[f];
		if field.FieldId == fieldId {
			return field;
		}

	}

	return Field{}
}
type Field struct {
	Id string `bson:"id" json:"id"`
	ParentId string `bson:"parentId" json:"parentId"`
	FieldId string `bson:"fieldId" json:"fieldId"`
	Name string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
	Placeholder string `bson:"placeholder" json:"placeholder"`
	IsDisabled bool `bson:"isDisabled" json:"isDisabled"`
	Type string `bson:"type" json:"type"`
	IsHidden bool `bson:"isHidden" json:"isHidden"`
	IsEditable bool `bson:"isEditable" json:"isEditable"`
	Items []Item `bson:"items" json:"items"`
	IsDirty bool `bson:"isDirty" json:"isDirty"`
	Errors []string `bson:"errors" json:"errors"`
	Validations []Validation `bson:"validations" json:"validations"`
	MinLength int `bson:"minLength" json:"minLength"`
	MaxLength int `bson:"maxLength" json:"maxLength"`
	IsMandatory bool `bson:"isMandatory" json:"isMandatory"`
	RegexValue string `bson:"regexValue" json:"regexValue"`
	Title string `bson:"title" json:"title"`
	SortOrder* string `bson:"sortOrder" json:"sortOrder"`
}

type Item struct {
	Id string `bson:"id" json:"id"`
	Item string `bson:"item" json:"item"`
	Values map[string]Item `json:"values"`
}

type Validation struct {
	Id string `bson:"id" json:"id"`
	ParentId string `bson:"parentId" json:"parentId"`
	ValidationId string `bson:"validationId" json:"validationId"`
	ValidationName string `bson:"validationName" json:"validationName"`
	Message string `bson:"message" json:"message"`
	IsValid bool 
	IsEnabled bool `bson:"isEnabled" json:"isEnabled"`
}