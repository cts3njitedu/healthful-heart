package utils

import (
	"github.com/cts3njitedu/healthful-heart/models"
)

func FindSection(sectionName string, page models.Page) (models.Section) {
	for s := range page.Sections {
		var section = page.Sections[s];
		if section.SectionId == sectionName {
			return section;
		}

	}

	return models.Section{};
}

func CloneSection(section models.Section) (models.Section) {
	newSection := models.Section{};
	newSection.Id = section.Id;
	newSection.ParentId = section.ParentId;
	newSection.SectionId = section.SectionId;
	newSection.Errors = append(section.Errors[:0:0], section.Errors...)
	newFields := make([]models.Field, 0, len(section.Fields))
	for _, field := range section.Fields {
		newFields = append(newFields, CloneField(field))
	}
	newSection.Fields = newFields;
	return newSection;

}

func CloneField(field models.Field) (models.Field) {
	newField  := models.Field{};
	newField.Id = field.Id;
	newField.FieldId = field.FieldId
	newField.ParentId = field.ParentId;
	newField.Name = field.Name;
	newField.Value = field.Value;
	newField.Placeholder = field.Placeholder;
	newField.IsDisabled = field.IsDisabled;
	newField.Type = field.Type;
	newField.IsHidden = field.IsHidden;
	newField.IsDirty = field.IsDirty;
	newField.Errors = append(newField.Errors[:0:0], newField.Errors...)
	newField.MinLength = field.MinLength;
	newField.MaxLength = field.MaxLength;
	newField.IsMandatory = field.IsMandatory;
	newField.RegexValue = field.RegexValue;
	newItems := make([]models.Item, 0, len(field.Items));
	for _, item := range field.Items {
		newItems = append(newItems, CloneItem(item))
	}
	newField.Items = newItems
	newValidations := make([]models.Validation, 0, len(field.Validations))
	for _, validation := range field.Validations {
		newValidations = append(newValidations, CloneValidation(validation))
	}
	newField.Validations = newValidations
	return newField;
}

func CloneItem(item models.Item) (models.Item) {
	newItem := models.Item{};
	newItem.Id = item.Id;
	newItem.Item = item.Item;
	return newItem;
}

func CloneValidation(validation models.Validation) (models.Validation) {
	newValidation := models.Validation{};
	newValidation.Id = validation.Id
	newValidation.ParentId = validation.ParentId;
	newValidation.ValidationId = validation.ValidationId;
	newValidation.ValidationName = validation.ValidationName;
	newValidation.Message = validation.Message;
	newValidation.IsValid = validation.IsValid;
	newValidation.IsEnabled = validation.IsEnabled;
	return newValidation;
}