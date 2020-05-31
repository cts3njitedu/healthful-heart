package utils

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	"strings"
	"strconv"
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

func SqlQueryBuilder(queryOptions models.QueryOptions, columns map[string]models.QueryOptions, sortMap map[string]models.QueryOptions, tableName string) (string, []interface{}) {
	length := 0;
	selectQuery := SelectQuery(queryOptions, columns);
	whereQuery, whereValues := WhereQuery(queryOptions)
	length += len(whereValues)	
	orderQuery := OrderQuery(queryOptions, columns, sortMap)
	inQuery := InQuery(queryOptions)
	notInQuery := NotInQuery(queryOptions)
	limitAndOffsetQuery := LimitAndOffsetQuery(queryOptions)
	fmt.Printf("Where Query Builder:%+v\n",whereQuery)
	values := make([]interface{}, 0, length);
	totalQuery := "SELECT " + strings.Join(selectQuery, " , ") + " FROM " + fmt.Sprintf("%v",tableName) + " WHERE " + strings.Join(whereQuery, " AND ") +
	strings.Join(inQuery, " ") + strings.Join(notInQuery, " ") + " ORDER BY " + fmt.Sprintf("%v", strings.Join(orderQuery, " , ")) +
	limitAndOffsetQuery;
	values = append(values, whereValues...)
	fmt.Println("Query:", totalQuery)
	fmt.Printf("Values: %+v\n",values)
	return totalQuery, values
}


func SelectQuery(queryOptions models.QueryOptions, columns map[string]models.QueryOptions) ([]string) {
	selectQuery := []string{"*"};
	if queryOptions.Select != nil && len(queryOptions.Select) > 0 {
		selectQuery = make([]string, 0 , len(queryOptions.Select));
		for _, v := range queryOptions.Select {
			if _, ok := columns[strings.ToLower(v)]; ok {
				selectQuery = append(selectQuery, strings.ToUpper(v));
			}
		}
		
	}
	fmt.Printf("Select Query: %+v\n", selectQuery)
	return selectQuery
}

func WhereQuery(queryOptions models.QueryOptions) ([]string, []interface{}) {
	whereQuery := []string{"1=?"};
	var whereValues []interface{};
	whereValues = append(whereValues, 1);
	whereEqual := queryOptions.WhereEqual;
	if queryOptions.Where != nil && len(queryOptions.Where) > 0 {
		whereQuery = make([]string, 0, len(queryOptions.Where))
		whereValues = make([]interface{}, 0, len(queryOptions.Where))
		for k, v := range queryOptions.Where {
			_, isWhere := whereEqual[k];
			if !queryOptions.IsEqual || !isWhere {
				if nv, ok := v.(string); ok {
					whereValues = append(whereValues, nv + "%")
					whereQuery = append(whereQuery, fmt.Sprintf("%s LIKE ?", k))
				} 
			} else {
				whereValues = append(whereValues, v)
				whereQuery = append(whereQuery, fmt.Sprintf("UPPER(%s) = UPPER(?)", k))
			}
		}
	}
	newWhereQuery := make([]string, len(whereQuery))
	copy(newWhereQuery, whereQuery)
	
	fmt.Printf("Where Query: %+v, Where Values: %+v\n", whereQuery, whereValues)
	return newWhereQuery, whereValues
}

func OrderQuery(queryOptions models.QueryOptions, columns map[string]models.QueryOptions, sortMap map[string]models.QueryOptions) ([]string) {
	orderQuery := []string {"1"}
	if queryOptions.Order != nil && len(queryOptions.Order) > 0 {
		orderQuery = make([]string, 0, len(queryOptions.Order))
		for k, v := range queryOptions.Order {
			if _, ok := columns[strings.ToLower(k)]; ok {
				if _, sok := sortMap[strings.ToLower(v)]; sok {
					orderQuery = append(orderQuery, k + " " + v)
				}
			}
		}
	}
	fmt.Printf("Order Query: %+v\n", orderQuery);
	return orderQuery
}

func InQuery(queryOptions models.QueryOptions) ([]string) {
	inQuery := []string{""}
	if queryOptions.In != nil && len(queryOptions.In) > 0 {
		for k, v := range queryOptions.In {
			inQuery = append(inQuery, fmt.Sprintf(" AND %s in (%s)", k, v))
		}
	}

	fmt.Printf("In Query: %+v\n", inQuery)
	return inQuery
}

func NotInQuery(queryOptions models.QueryOptions) ([]string) {
	notInQuery := []string{""}
	if queryOptions.NotIn != nil && len(queryOptions.NotIn) > 0 {
		for k, v := range queryOptions.NotIn {
			notInQuery = append(notInQuery, fmt.Sprintf(" AND %s NOT IN (%s)", k, v))
		}
	}
	fmt.Printf("Not In Query: %+v\n", notInQuery)
	return notInQuery
}

func LimitAndOffsetQuery(queryOptions models.QueryOptions) (string) {
	limitAndOffsetQuery := "";
	if queryOptions.Limit !=0 || queryOptions.Offset !=0 {
		limitAndOffsetQuery = " LIMIT " + strconv.FormatInt(queryOptions.Limit, 10) + " OFFSET " + strconv.FormatInt(queryOptions.Offset, 10) 
	}

	fmt.Printf("Limit and Filter Offset: %+v\n", limitAndOffsetQuery);
	return limitAndOffsetQuery
}