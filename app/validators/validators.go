package validators

import (
	"regexp"
	"strconv"
)

func IsEmpty(s string) bool {
	return (s == "")
}

func IsLongerThan(s string, min int) bool {
	return (len(s) > min)
}

func IsLongerOrEqualThan(s string, min int) bool {
	return (len(s) >= min)
}

func IsShorterThan(s string, max int) bool {
	return (len(s) < max)
}

func IsShorterOrEqualThan(s string, max int) bool {
	return (len(s) <= max)
}

func IsNumeric(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func HasNumber(s string) bool {
	for _, value := range s {
		switch {
		case value >= '0' && value <= '9':
			return true
		}
	}
	return false
}

func IsEmail(s string) bool {
	var pattern string = `[^@]+@[^@]+\.[^@]+`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

func IsURL(s string) bool {
	var pattern string = `((https?):((//)|(\\\\))+([\w\d:#@%/;$()~_?\+-=\\\.&](#!)?)*)`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}
