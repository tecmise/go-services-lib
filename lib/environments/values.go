package environments

import (
	"os"
	"strconv"
)

func GetStringValue(param string) string {
	val := os.Getenv(param)
	if val == "" {
		panic("environment variable " + param + " is not set")
	}
	return val
}

func getIntValue(param string) *int {
	val := os.Getenv(param)
	if val == "" {
		return nil
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic("environment variable " + param + " is not a valid integer")
	}
	return &intVal
}

func GetIntValue(param string) int {
	val := getIntValue(param)
	if val == nil {
		panic("environment variable " + param + " is not set")
	}
	return *val
}

func GetStringValueDefault(param string, _default string) string {
	val := GetStringValue(param)
	if val == "" {
		return _default
	}
	return val
}

func GetIntValueDefault(param string, _default int) int {
	val := getIntValue(param)
	if val == nil {
		return _default
	}
	return *val
}
