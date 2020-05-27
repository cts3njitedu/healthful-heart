package services

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"strings"
	"strconv"
	"fmt"
)

type GroupParserService struct{

}


func NewGroupParserService() *GroupParserService {
	return &GroupParserService{}
}


func (parser *GroupParserService) GetGroups(workoutType string, groupText string, categoryCode string) ([]models.Group) {
	groups := strings.Split(groupText, ",");
	fmt.Printf("Groups for work type %s are %s\n", workoutType, groupText)
	var variation string
	parsedGroups := make([]models.Group, 0, len(groups));
	for i := range groups {
		group := strings.TrimSpace(groups[i])
		groupInfos := strings.Split(group, "x");
		var sets int
		var repetitions int
		var weight float32
		var weightOrRepetition string
		var parsedGroup models.Group
		lastInfo := groupInfos[len(groupInfos)-1]
		parenIndex := strings.Index(lastInfo, "(");
		if parenIndex >= 0 {
			weightOrRepetition = GetSubstring(lastInfo, 0, parenIndex);
			weightOrRepetition = strings.TrimSpace(weightOrRepetition)
			rightParen := strings.Index(lastInfo, ")")
			variation = GetSubstring(lastInfo, parenIndex+1, rightParen);
		} else {
			weightOrRepetition = lastInfo
		}
		if ((workoutType == models.PULLUPS && categoryCode == models.BACK) ||
			(workoutType == models.DIPS && categoryCode == models.CHEST) ||
			(workoutType == models.BICYCLE && categoryCode == models.ABS) ||
			(workoutType == models.LEG_RAISES && categoryCode == models.ABS) ||
			(workoutType == models.SITUPS && categoryCode == models.ABS) ||
			(workoutType == models.INCLINE_ABS && categoryCode == models.ABS)) {
			fmt.Printf("Workout type: %s, Category Code: %s\n", workoutType, categoryCode)
			s, err := strconv.Atoi(strings.TrimSpace(groupInfos[0]));
			if err != nil {
				fmt.Println(err)
				continue; 
			}
			sets = s
			reps, err := strconv.ParseInt(weightOrRepetition,10, 32)
			if err != nil {
				fmt.Println(err)
				continue; 
			}
			repetitions = int(reps)
		} else {
			// this means user only did 1 set
			if len(groupInfos) == 2 {
				sets = 1
				reps, err := strconv.Atoi(strings.TrimSpace(groupInfos[0]));
				if err != nil {
					fmt.Println(err)
					continue; 
				}
				repetitions = reps
				w, err := strconv.ParseFloat(weightOrRepetition,32)
				if err != nil {
					fmt.Println(err)
					continue; 
				}
				weight = float32(w)
			} else {
				s, err := strconv.Atoi(strings.TrimSpace(groupInfos[0]));
				if err != nil {
					fmt.Println(err)
					continue; 
				}
				sets = s
				reps, err := strconv.Atoi(strings.TrimSpace(groupInfos[1]));
				if err != nil {
					fmt.Println(err)
					continue; 
				}
				repetitions = reps
				w, err := strconv.ParseFloat(weightOrRepetition,32)
				if err != nil {
					fmt.Println(err)
					continue; 
				}
				weight = float32(w)

			}
		}
		parsedGroup.Sets = &sets
		parsedGroup.Repetitions = &repetitions
		parsedGroup.Weight = &weight
		parsedGroups = append(parsedGroups, parsedGroup)
	}

	for g := range parsedGroups {
		pGroup := &parsedGroups[g]
		pGroup.Variation = &variation
		sequence := int64(g+1)
		pGroup.Sequence = &sequence
	}
	return parsedGroups
}

func GetSubstring(s string, start int, end int) (string) {
	runes := []rune(s)
	substring := string(runes[start:end])
	return substring
}