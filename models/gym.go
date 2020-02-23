package models


type WorkoutDay struct {
	Workout_Day_Id int64
	User_Id int64
	Location_Id int64
	Workout_Date string
	Cre_Ts string
	Mod_Ts string
	Del_Ts string
	Workouts []Workout
	Location Location

}

type Workout struct {
	Workout_Id int64
	Workout_Day_Id int64
	Workout_Type_Cd string
	Cre_Ts string
	Mod_Ts string
	Del_Ts string
	Groups []Group
}

type Group struct {
	Group_Id int64
	Workout_Id int64
	Sequence int64
	Sets int
	Repetitions int
	Weight float32
	Variation string
	Cre_Ts string
	Mod_Ts string
	Del_Ts string
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

