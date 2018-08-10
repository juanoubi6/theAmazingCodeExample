package common

import (
	"strconv"
)

func StringToUint(stringToConvert string) (uint, error) {

	uintValue64, parseErr := strconv.ParseUint(stringToConvert, 10, 64)
	if parseErr != nil {
		return 0, parseErr
	}

	return uint(uintValue64), nil

}

func StringToInt(stringToConvert string) (int, error) {

	intValue, parseErr := strconv.Atoi(stringToConvert)
	if parseErr != nil {
		return 0, parseErr
	}

	return intValue, nil

}
