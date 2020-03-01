package models


type WorkoutDay struct {
	Workout_Day_Id int64 `gorm:"column:WORKOUT_DAY_ID;primary_key;not null; auto_increment"`
	User_Id int64 `gorm:"column:USER_ID"`
	Location_Id int64 `gorm:"column:LOCATION_ID"`
	Workout_Date string `gorm:"column:WORKOUT_DATE"`
	Cre_Ts *string `gorm:"column:CRE_TS"`
	Mod_Ts *string `gorm:"column:MOD_TS"`
	Del_Ts *string `gorm:"column:DEL_TS"`
	Workouts []Workout
	Location Location

}

type Workout struct {
	Workout_Id int64 `gorm:"column:WORKOUT_ID"`
	Workout_Day_Id int64 `gorm:"column:WORKOUT_DAY_ID"`
	Workout_Type_Cd string `gorm:"column:WORKOUT_TYPE_CD"`
	Workout_Name string `gorm:"-"`
	Cre_Ts *string `gorm:"column:CRE_TS"`
	Mod_Ts *string `gorm:"column:MOD_TS"`
	Del_Ts *string `gorm:"column:DEL_TS"`
	Groups []Group `gorm:"-"`
}

type Group struct {
	Group_Id int64 `gorm:"column:GROUP_ID"`
	Workout_Id int64 `gorm:"column:WORKOUT_ID"`
	Workout_Day_Id int64 `gorm:"column:WORKOUT_DAY_ID"`
	Sequence int64 `gorm:"column:SEQUENCE"`
	Sets int `gorm:"column:SETS"`
	Repetitions int `gorm:"column:REPETITIONS"`
	Weight float32 `gorm:"column:WEIGHT"`
	Variation string `gorm:"column:VARIATION"`
	Cre_Ts *string `gorm:"column:CRE_TS"`
	Mod_Ts *string `gorm:"column:MOD_TS"`
	Del_Ts *string `gorm:"column:DEL_TS"`
}

type Location struct {
	Location_Id int64
	Name string
	State string
	City string
	Country string
	Zipcode string
	Location string
}

