package services


import (
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	"errors"
	"strings"
	// "sort"
)
type WorkoutTypeService struct {
	workoutTypeRepository mysqlrepo.IWorkoutTypeRepository
	categoryRepository mysqlrepo.ICategoryRepository
}

var catCodeMap =  make(map[string]map[string]models.WorkoutType)
var catCodeTypeMap = make(map[string]map[string]string)

var categoryCdToName = make(map[string]string)

var categoryNameToCd = make(map[string]string)

var sortedCatAndWorkType = make([]models.SortedCategoryWorkoutType, 0 , 100)

func NewWorkoutTypeService(workoutTypeRepository mysqlrepo.IWorkoutTypeRepository,
	categoryRepository mysqlrepo.ICategoryRepository) *WorkoutTypeService {
	workTypes, err := workoutTypeRepository.GetWorkoutTypes();

	if err != nil {
		fmt.Println("Unable to pre load workout types")
	}
	sortedWorkoutsMap := make(map[string][]models.WorkoutType)
	for _,workType := range workTypes {
		sortedWorkoutsMap[workType.Category_Cd] = append(sortedWorkoutsMap[workType.Category_Cd], workType)
		if catCodeMap[workType.Category_Cd] == nil {
			workTypeNameMap := make(map[string]models.WorkoutType);
			workTypeName := strings.ToUpper(workType.Name);
			workTypeNameMap[workTypeName] = workType;
			catCodeMap[workType.Category_Cd]=workTypeNameMap
		} else {
			workTypeNameMap := catCodeMap[workType.Category_Cd]
			workTypeName := strings.ToUpper(workType.Name);
			workTypeNameMap[workTypeName] = workType;
		}
		if catCodeTypeMap[workType.Category_Cd] == nil {
			workTypeMap := make(map[string]string)
			workTypeMap[workType.Workout_Type_Cd] = workType.Name
			catCodeTypeMap[workType.Category_Cd] = workTypeMap
		} else {
			workTypeMap := catCodeTypeMap[workType.Category_Cd]
			workTypeMap[workType.Workout_Type_Cd] = workType.Name
		}
	} 


	categories, err := categoryRepository.GetCategories();

	if err != nil {
		fmt.Println("Unable to pre load categories")
	}
	fmt.Printf("Categories: %+v\n", categories)
	for _,category := range categories {
		catWorkSort := models.SortedCategoryWorkoutType{}
		catWorkSort.Category_Cd = category.Category_Cd
		catWorkSort.Category_Name = category.Category_Name
		catWorkSort.WorkoutTypes = sortedWorkoutsMap[category.Category_Cd]
		sortedCatAndWorkType = append(sortedCatAndWorkType, catWorkSort)
		categoryCdToName[category.Category_Cd] = category.Category_Name
		categoryNameToCd[category.Category_Name] = category.Category_Cd
	}
	return &WorkoutTypeService{workoutTypeRepository, categoryRepository}
}


func (serv *WorkoutTypeService) GetWorkoutTypeCode(categoryCd string, workoutTypeName string) string {
	fmt.Printf("Category Code: %v, Workout Type Name: %v\n", categoryCd, workoutTypeName)
	if catCodeMap[categoryCd] == nil {
		fmt.Println("Couldn't find category code in map")
		return getSubstring(workoutTypeName, 0, 5);
	} else {
		workTypeNameMap:=catCodeMap[categoryCd];
		workoutTypeName = strings.ToUpper(workoutTypeName)
		if workTypeNameMap[workoutTypeName] == (models.WorkoutType{}) {
			fmt.Println("Couldn't find work type code for:", workoutTypeName)
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

func (serv *WorkoutTypeService) GetCategories() (map[string]string, map[string]string) {
	newNameToCd := make(map[string]string)
	for k, v := range categoryNameToCd {
		newNameToCd[k] = v
	}
	newCdToName := make(map[string]string)
	for k,v := range categoryCdToName {
		newCdToName[k] = v
	}
	return newNameToCd, newCdToName
}

func (serv *WorkoutTypeService) GetCategoriesAndWorkouts() (map[string]map[string]models.WorkoutType) {
	newCatWorkoutMap := make(map[string]map[string]models.WorkoutType)

	for c, wm := range catCodeMap {
		newWorkoutMap := make(map[string]models.WorkoutType)
		for wc, wt := range wm {
			newWorkoutMap[wc] = models.WorkoutType{
				Name: wt.Name,
				Category_Cd: wt.Category_Cd,
				Workout_Type_Cd: wt.Workout_Type_Cd,
			}
		}
		newCatWorkoutMap[c] = newWorkoutMap
	}
	return newCatWorkoutMap
}

func (serv *WorkoutTypeService) GetCategoriesAndWorkoutTypes() (map[string]map[string]string) {
	newCatWorkoutMap := make(map[string]map[string]string)

	for c, wm := range catCodeTypeMap {
		newWorkoutMap := make(map[string]string)
		for wc, wt := range wm {
			newWorkoutMap[wc] = wt
		}
		newCatWorkoutMap[c] = newWorkoutMap
	}
	return newCatWorkoutMap
}

func (serv *WorkoutTypeService) GetSortedCategoriesAndWorkoutTypes() ([]models.SortedCategoryWorkoutType) {
	newSorts := make([]models.SortedCategoryWorkoutType, 0, len(sortedCatAndWorkType))
	for _, s := range sortedCatAndWorkType {
		newSort := models.SortedCategoryWorkoutType{}
		newSort.Category_Cd = s.Category_Cd;
		newSort.Category_Name = s.Category_Name;
		newTypes := make([]models.WorkoutType, 0, len(s.WorkoutTypes))
		for _, wok := range s.WorkoutTypes {
			newType := models.WorkoutType{}
			newType.Workout_Type_Cd = wok.Workout_Type_Cd;
			newType.Name = wok.Name;
			newType.Category_Cd = wok.Category_Cd;
			newTypes = append(newTypes, newType)
		}
		newSort.WorkoutTypes = newTypes;
		newSorts = append(newSorts, newSort)
	}
	return newSorts
}
func getSubstring(s string, start int, end int) (string) {
	runes := []rune(s)
	substring := string(runes[start:end])
	return substring
}