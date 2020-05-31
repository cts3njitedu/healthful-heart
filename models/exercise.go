package models

type WorkoutType struct {
	Workout_Type_Desc string `gorm:"column:WORKOUT_TYPE_DESC"`
	Category_Cd string `gorm:"column:CATEGORY_CD"`
	Workout_Type_Id int64 `gorm:"column:WORKOUT_TYPE_ID"`
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
	PULLUPS = "Pull ups"
	DIPS = "Dips"
	BICYCLE = "Bicycle"
	LEG_RAISES = "Leg Raises"
	SITUPS = "Situps"
	INCLINE_ABS = "Incline Abs"
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