package models

type WorkoutType struct {
	Name string `gorm:"column:NAME"`
	Category_Cd string `gorm:"column:CATEGORY_CD"`
	Workout_Type_Cd string `gorm:"column:WORKOUT_TYPE_CD"`
}

type Category struct {
	Category_Cd string `gorm:"column:CATEGORY_CD"`
	Category_Name string `gorm:"column:CATEGORY_NAME"`
}

type SortedCategoryWorkoutType struct {
	Category
	WorkoutTypes []WorkoutType

}

const (
	BACK = "BK"
	ABS = "AB"
	CHEST = "CH"
	TRICEPS = "TR"
	LEGS = "LG"
	BICEPS = "BC"
	SHOULDERS = "SH"
	DEFAULT = "DF"
	PULLUPS = "BK16"
	DIPS = "CH5"
	BICYCLE = "AB3"
	LEG_RAISES = "AB5"
	SITUPS = "AB6"
	INCLINE_ABS = "AB4"
)


func GetCategoryCodeFromName(catName string) (string) {
	switch catName {
	case "Back":
		return "BK"
	case "Abs":
		return "AB"
	case "Chest":
		return "CH"
	case "Triceps":
		return "TR"
	case "Legs":
		return "LG"
	case "Biceps":
		return "BC"
	case "Shoulders":
		return "SH"
	default:
		return "DF"
	}
}