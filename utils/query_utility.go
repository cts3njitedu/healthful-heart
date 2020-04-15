package utils

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
)

func QueryBuildSort(sort map[string]models.HeartSort, section models.Section) (map[string]string) {
	if sort == nil {
		return nil
	}
	sortMap := make(map[string]string);
	for f := range section.Fields {
		var field  = &section.Fields[f]
		if value, ok :=  sort[field.Name]; ok {
			sortMap[field.Name] = value.SortOrder;
		}
		
	}
	fmt.Printf("Sort slice: %+v\n", sortMap)
	
	return sortMap
}

// func SqlQueryBuilder()