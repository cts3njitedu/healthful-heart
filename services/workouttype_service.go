package services


import (
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	"errors"
)
type WorkoutTypeService struct {
	workoutTypeRepository mysqlrepo.IWorkoutTypeRepository
	categoryRepository mysqlrepo.ICategoryRepository
}

var catCodeMap =  make(map[string]map[string]models.WorkoutType)

var categoryCdToName = make(map[string]string)

var categoryNameToCd = make(map[string]string)


func NewWorkoutTypeService(workoutTypeRepository mysqlrepo.IWorkoutTypeRepository,
	categoryRepository mysqlrepo.ICategoryRepository) *WorkoutTypeService {
	workTypes, err := workoutTypeRepository.GetWorkoutTypes();

	if err != nil {
		fmt.Println("Unable to pre load workout types")
	}

	for _,workType := range workTypes {
		if catCodeMap[workType.Category_Cd] == nil {
			workTypeNameMap := make(map[string]models.WorkoutType);
			workTypeNameMap[workType.Name] = workType;
			catCodeMap[workType.Category_Cd]=workTypeNameMap
		} else {
			workTypeNameMap := catCodeMap[workType.Category_Cd]
			workTypeNameMap[workType.Name] = workType;
		}
	} 

	categories, err := categoryRepository.GetCategories();

	if err != nil {
		fmt.Println("Unable to pre load categories")
	}
	fmt.Printf("Categories: %+v\n", categories)
	// for _,category := range categories {
	// 	categoryCdToName[category.Category_Cd] = category.Category_Name
	// 	categoryNameToCd[category.Category_Name] = category.Category_Cd
	// }
	return &WorkoutTypeService{workoutTypeRepository, categoryRepository}
}


func (serv *WorkoutTypeService) GetWorkoutTypeCode(categoryCd string, workoutTypeName string) string {
	if catCodeMap[categoryCd] == nil {
		fmt.Println("Couldn't find category code in map")
		return getSubstring(workoutTypeName, 0, 5);
	} else {
		workTypeNameMap:=catCodeMap[categoryCd];
		if workTypeNameMap[workoutTypeName] == (models.WorkoutType{}) {
			return getSubstring(workoutTypeName, 0, 5);
		}
		return workTypeNameMap[workoutTypeName].Workout_Type_Cd
	}
}

func (serv *WorkoutTypeService) GetCategoryNameFromCode(categoryCd string) (string, error) {
	if categoryCdToName[categoryCd] == "" {
		return "", errors.New("No name for code")
	} 
	return categoryCdToName[categoryCd], nil
}

func (serv *WorkoutTypeService) GetCategoryCodeFromName(categoryName string) (string, error) {
	if categoryNameToCd[categoryName] == "" {
		return "", errors.New("No code for name")
	} 
	return categoryNameToCd[categoryName], nil
}
func getSubstring(s string, start int, end int) (string) {
	runes := []rune(s)
	substring := string(runes[start:end])
	return substring
}