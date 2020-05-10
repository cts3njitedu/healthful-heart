package models


type WorkoutDay struct {
	Workout_Day_Id int64 `gorm:"column:WORKOUT_DAY_ID;primary_key;not null; auto_increment" json:"workoutDayId"`
	User_Id int64 `gorm:"column:USER_ID" json:"userId"`
	Location_Id int64 `gorm:"column:LOCATION_ID" json:"-"`
	Workout_Date string `gorm:"column:WORKOUT_DATE" json:"-"`
	Month string `gorm:"-" json:"month"`
	MonthId int `gorm:"-" json:"monthId"`
	Year int `gorm:"-" json:"year"`
	NumberOfDays int `gorm:"-" json:"numberOfDays"`
	Day int `gorm:"-" json:"day"`
	Cre_Ts *string `gorm:"column:CRE_TS" json:"-"`
	Mod_Ts *string `gorm:"column:MOD_TS" json:"-"`
	Del_Ts *string `gorm:"column:DEL_TS" json:"-"`
	Workouts []Workout `gorm:"-" json:"-"`
	Location Location `gorm:"-" json:"location"`
	Version_Nb int64 `gorm:"column:VERSION_NB" json:"versionNb"`
	Workout_File_Id* int64 `gorm:"column:WORKOUT_FILE_ID" json:"fileId"`

}

type Workout struct {
	Workout_Id int64 `gorm:"column:WORKOUT_ID;primary_key;not null; auto_increment"`
	Workout_Day_Id int64 `gorm:"column:WORKOUT_DAY_ID"`
	Workout_Type_Cd string `gorm:"column:WORKOUT_TYPE_CD"`
	Workout_Name string `gorm:"-"`
	Cre_Ts *string `gorm:"column:CRE_TS"`
	Mod_Ts *string `gorm:"column:MOD_TS"`
	Del_Ts *string `gorm:"column:DEL_TS"`
	Groups []Group `gorm:"-"`
	Version_Nb int64 `gorm:"column:VERSION_NB" json:"versionNb"`
}

type Group struct {
	Group_Id int64 `gorm:"column:GROUP_ID;primary_key;not null; auto_increment"`
	Workout_Id int64 `gorm:"column:WORKOUT_ID"`
	Workout_Day_Id int64 `gorm:"column:WORKOUT_DAY_ID"`
	Sequence int64 `gorm:"column:SEQUENCE"`
	Sets int `gorm:"column:SETS"`
	Repetitions int `gorm:"column:REPETITIONS"`
	Weight float32 `gorm:"column:WEIGHT"`
	Duration float32 `gorm:"column:DURATION"`
	Variation string `gorm:"column:VARIATION"`
	Cre_Ts *string `gorm:"column:CRE_TS"`
	Mod_Ts *string `gorm:"column:MOD_TS"`
	Del_Ts *string `gorm:"column:DEL_TS"`
	Version_Nb int64 `gorm:"column:VERSION_NB" json:"versionNb"`
}

type Location struct {
	Location_Id int64 `gorm:"column:LOCATION_ID" json:"locationId"`
	Name string `gorm:"column:NAME" json:"gymName"`
	State string `gorm:"column:STATE" json:"state"`
	City string `gorm:"column:CITY" json:"city"`
	Country string `gorm:"column:COUNTRY" json:"country"`
	Zipcode string `gorm:"column:ZIPCODE" json:"zipCode"`
	Location string `gorm:"column:LOCATION" json:"locationName"`
}

