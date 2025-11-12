package vars

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func GetStringValue(param string) string {
	if val := os.Getenv(param); val != "" {
		return val
	}
	if v, ok := readDotEnvVar(param); ok {
		return v
	}
	panic("environment variable " + param + " is not set")
}

func getIntValue(param string) *int {
	val := os.Getenv(param)
	if val == "" {
		if v, ok := readDotEnvVar(param); ok {
			val = v
		} else {
			return nil
		}
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
	if val := os.Getenv(param); val != "" {
		return val
	}
	if v, ok := readDotEnvVar(param); ok {
		return v
	}
	return _default
}

func GetIntValueDefault(param string, _default int) int {
	val := getIntValue(param)
	if val == nil {
		return _default
	}
	return *val
}

func getFloat64Value(param string) *float64 {
	val := os.Getenv(param)
	if val == "" {
		if v, ok := readDotEnvVar(param); ok {
			val = v
		} else {
			return nil
		}
	}
	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic("environment variable " + param + " is not a valid float64")
	}
	return &floatVal
}

func GetFloat64Value(param string) float64 {
	val := getFloat64Value(param)
	if val == nil {
		panic("environment variable " + param + " is not set")
	}
	return *val
}

func GetFloat64ValueDefault(param string, _default float64) float64 {
	val := getFloat64Value(param)
	if val == nil {
		return _default
	}
	return *val
}

func IsDebugEnabled() bool {
	if val := os.Getenv("DEBUG"); val != "" {
		return strings.ToLower(strings.TrimSpace(val)) == "true"
	}
	if v, ok := readDotEnvVar("DEBUG"); ok {
		return strings.ToLower(strings.TrimSpace(v)) == "true"
	}
	return false
}

func IsProduction() bool {
	if val := os.Getenv("ENVIRONMENT"); val != "" {
		return strings.ToLower(strings.TrimSpace(val)) == "production"
	}
	if v, ok := readDotEnvVar("ENVIRONMENT"); ok {
		return strings.ToLower(strings.TrimSpace(v)) == "production"
	}
	return false
}

func readDotEnvVar(param string) (string, bool) {
	f, err := os.Open(".env")
	if err != nil {
		return "", false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if key != param {
			continue
		}
		val := strings.TrimSpace(parts[1])
		if len(val) >= 2 {
			if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) ||
				(strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
				val = val[1 : len(val)-1]
			}
		}
		return val, true
	}
	return "", false
}
