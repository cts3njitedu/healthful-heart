package models


func GetCategoryCodeFromName(catName string) (string) {
	switch catName {
	case "Back":
		return "BACK"
	case "Abs":
		return "ABS"
	case "Chest":
		return "CHEST"
	case "Triceps":
		return "TRICEPS"
	case "Legs":
		return "LEGS"
	case "Biceps":
		return "BICEPS"
	case "Shoulders":
		return "SHLDRS"
	default:
		return "DEFAULT"
	}
}