package models

type WorkoutFile struct {
	Workout_File_Id int64 `json:"fileId"`
	File_Name string `json:"fileName"`
	File_Path string `json:"filePath"`
	Status string `json:"status"`
	User_Id int64 `json:"userId"`
	Version_Nb int64 `json:"versionNb"`
	Location_Id int64 `json:"locationId"`
	Cre_Ts *string
	Mod_Ts *string
}