package runner

import (
	"fmt"
	"regexp"
)

func CheckEn(myString string) bool {
	b, err := regexp.MatchString("^[a-zA-Z\\s.'-]+$", myString)
	if err != nil {
		fmt.Println("Check SearchKeyword Error!")
	}
	return b
}
