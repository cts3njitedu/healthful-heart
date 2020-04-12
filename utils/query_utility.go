package utils

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
)

func QueryBuildSort(sort map[string]models.HeartSort, section models.Section) ([]string) {
	if sort == nil {
		return nil
	}
	sortSlice := make([]string, 0, len(section.Fields))
	for f := range section.Fields {
		var field  = &section.Fields[f]
		if value, ok :=  sort[field.Name]; ok {
			sortString := field.Name + " " + value.SortOrder;
			sortSlice = append(sortSlice, sortString);
		}
		
	}
	fmt.Printf("Sort slice: %+v\n", sortSlice)
	
	return sortSlice
}