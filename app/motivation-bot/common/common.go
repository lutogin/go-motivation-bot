package common

import (
	"regexp"
)

func CriticErrorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func Is24TimeFormat(inpTime string) bool {
	// regular expression for matching HH:MM format
	re := regexp.MustCompile(`^([01][0-9]|2[0-3]):[0|3]0$`)

	return re.MatchString(inpTime)
}
