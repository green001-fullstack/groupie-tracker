package utils

import (
	"strings"
)


func FormattedDatesAndLocation(str string) string {
	newString := strings.Split(str, "-")
	var eachStringSlice []string

	for i := 0; i < len(newString); i++ {
		newString[i] = strings.ReplaceAll(newString[i], "_", " ")
		eachStringSlice = strings.Split(newString[i], " ")
		var newEachStringSlice []string
		for j := 0; j < len(eachStringSlice); j++ {
			if len(eachStringSlice[j])> 0{
				eachStringSlice[j] = strings.ToUpper(eachStringSlice[j][:1]) + strings.ToLower(eachStringSlice[j][1:])
			}
			newEachStringSlice = append(newEachStringSlice, eachStringSlice[j])
		}
		newString[i] = strings.Join(newEachStringSlice, " ")
	}
	result := strings.Join(newString, ", ")
	return result

}

