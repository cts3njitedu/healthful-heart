package services


import (
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	"errors"
	"strings"
	"strconv"
	// "sort"
)
type WorkoutTypeService struct {
	workoutTypeRepository mysqlrepo.IWorkoutTypeRepository
	categoryRepository mysqlrepo.ICategoryRepository
}


func NewWorkoutTypeService(workoutTypeRepository mysqlrepo.IWorkoutTypeRepository,
	categoryRepository mysqlrepo.ICategoryRepository) *WorkoutTypeService {
	return &WorkoutTypeService{workoutTypeRepository, categoryRepository}
}


func (serv *WorkoutTypeService) GetWorkoutType(categoryCd string, workoutTypeName string) (models.WorkoutType, error) {
	fmt.Printf("Category Code: %v, Workout Type Name: %v\n", categoryCd, workoutTypeName)
	queryOptions := models.QueryOptions{}
	queryOptions.Where = map[string]interface{} {
		"category_cd" : categoryCd,
		"workout_type_desc" : workoutTypeName,
	}

	queryOptions.WhereEqual = map[string]bool {
		"category_cd" : true,
		"workout_type_desc" : true,
	}
	queryOptions.IsEqual = true;
	workoutTypes, err := serv.workoutTypeRepository.GetWorkoutTypes(queryOptions)

	if len(workoutTypes) != 1 || err != nil {
		return models.WorkoutType{}, err
	} 
	return workoutTypes[0], nil

}

func (serv *WorkoutTypeService) GetWorkoutTypeByIds(ids []int64) (map[int64]models.WorkoutType) {
	queryOptions := models.QueryOptions{}
	idStrings := make([]string, 0, len(ids))
	for _, v := range ids {
		s := strconv.FormatInt(v, 10)
		idStrings = append(idStrings, s)
	}
	inQuery := strings.Join(idStrings, ",");
	queryOptions.In = map[string]string {
		"workout_type_id" : inQuery,
	}
	workoutTypes, err := serv.workoutTypeRepository.GetWorkoutTypes(queryOptions);
	if err != nil {
		fmt.Println("Something went wrong getting workout type ids")
	}
	workMap := make(map[int64]models.WorkoutType)
	for _, w := range workoutTypes {
		workMap[w.Workout_Type_Id] = w
	}
	return workMap
} 


func (serv *WorkoutTypeService) GetCategoryNameFromCode(categoryCd string) (string, error) {
	queryOptions := models.QueryOptions{}
	queryOptions.Where = map[string]interface{} {
		"category_cd" : categoryCd,
	}
	queryOptions.WhereEqual = map[string]bool {
		"category_cd" : true,
	}
	queryOptions.IsEqual = true;
	
	categories, err := serv.categoryRepository.GetCategoriesByParams(queryOptions)
	if len(categories) == 0 || err != nil {
		return "", errors.New("No Name For Code")
	}
	return categories[0].Category_Name, nil
}

func (serv *WorkoutTypeService) GetCategoryCodeFromName(categoryName string) (string, error) {
	queryOptions := models.QueryOptions{}
	queryOptions.Where = map[string]interface{} {
		"category_name" : categoryName,
	}
	queryOptions.WhereEqual = map[string]bool {
		"category_name" : true,
	}
	queryOptions.IsEqual = true;
	
	categories, err := serv.categoryRepository.GetCategoriesByParams(queryOptions)
	if len(categories) == 0 || err != nil {
		return "", errors.New("No Code For Category Name")
	}
	return categories[0].Category_Cd, nil
}

func (serv *WorkoutTypeService) GetCategories() (map[string]string, map[string]string) {
	queryOptions := models.QueryOptions{}
	categories, err := serv.categoryRepository.GetCategoriesByParams(queryOptions)
	newNameToCd := make(map[string]string)
	newCdToName := make(map[string]string)
	if err != nil {
		fmt.Println("Something went wrong getting categories")
	}
	for _, category := range categories {
		newNameToCd[category.Category_Name] = category.Category_Cd;
		newCdToName[category.Category_Cd] = category.Category_Name
	}
	
	return newNameToCd, newCdToName
}

func (serv *WorkoutTypeService) GetCategoriesAndWorkoutsMap(catCode string) (map[string]map[int64]models.WorkoutType) {
	newCatWorkoutMap := make(map[string]map[int64]models.WorkoutType)
	queryOptions := models.QueryOptions{}
	queryOptions.Order = map[string]string {
		"workout_type_desc": "asc",
	}
	if len(catCode) != 0 {
		queryOptions.Where = map[string]interface{} {
			"category_cd": catCode,
		}
		queryOptions.WhereEqual = map[string]bool {
			"category_cd": true,
		}
		queryOptions.IsEqual = true
	}
	workTypes, err := serv.workoutTypeRepository.GetWorkoutTypes(queryOptions)
	if err != nil {
		fmt.Println("Something went wrong getting category and worktypes map")
	}
	for _, wk := range workTypes {
		if cat, ok := newCatWorkoutMap[wk.Category_Cd]; ok {
			cat[wk.Workout_Type_Id] = wk
		} else {
			newCat := make(map[int64]models.WorkoutType)
			newCat[wk.Workout_Type_Id] = wk
			newCatWorkoutMap[wk.Category_Cd] = newCat
		}
	}
	return newCatWorkoutMap
}


func (serv *WorkoutTypeService) GetSortedCategoriesAndWorkoutTypes() ([]models.SortedCategoryWorkoutType) {
	
	queryOptions := models.QueryOptions{}
	queryOptions.Order = map[string]string{
		"category_name" : "asc",
	}
	categories, err := serv.categoryRepository.GetCategoriesByParams(queryOptions)
	if err != nil {
		fmt.Println("Something went wrong getting categories")
	}
	sortCats := make([]models.SortedCategoryWorkoutType, 0, 100)

	for c := range categories {
		cat := categories[c]
		sortCat := models.SortedCategoryWorkoutType{}
		sortCat.Category_Cd = cat.Category_Cd;
		sortCat.Category_Name = cat.Category_Name
		workOptions := models.QueryOptions{}
		workOptions.Order = map[string]string {
			"workout_type_desc" : "asc",
		}
		workOptions.Where = map[string]interface{} {
			"category_cd" : cat.Category_Cd,
		}
		workOptions.WhereEqual = map[string]bool {
			"category_cd" : true,
		}
		workOptions.IsEqual = true;
		workoutTypes, err := serv.workoutTypeRepository.GetWorkoutTypes(workOptions)
		if err != nil {
			fmt.Printf("Error Getting Workout Types: %+v", cat.Category_Cd)
		}
		// fmt.Printf("Workout Types Chicken: %+v\n", workoutTypes)
		sortCat.WorkoutTypes = workoutTypes
		sortCats = append(sortCats, sortCat)
	}

	
	return sortCats
}