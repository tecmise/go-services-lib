package vars

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type (
	Configuration interface {
		GetStringValue(param string) string
		GetIntValue(param string) int
		GetStringValueDefault(param string, _default string) string
		GetIntValueDefault(param string, _default int) int
		GetFloat64Value(param string) float64
		GetFloat64ValueDefault(param string, _default float64) float64
		IsDebugEnabled() bool
		IsProduction() bool
	}

	configuration struct {
		repository string
		context    string
	}
)

func NewConfiguration(repository string, context string) Configuration {
	return &configuration{
		repository: repository,
		context:    context,
	}
}

func (c *configuration) GetStringValue(param string) string {
	if v, ok := c.readSecretStoreVar(param); ok {
		return v
	}
	if val := os.Getenv(param); val != "" {
		return val
	}
	if v, ok := readDotEnvVar(param); ok {
		return v
	}
	panic("environment variable " + param + " is not set")
}

func (c *configuration) getIntValue(param string) *int {
	var val string
	if v, ok := c.readSecretStoreVar(param); ok {
		val = v
	} else if v := os.Getenv(param); v != "" {
		val = v
	} else if v, ok := readDotEnvVar(param); ok {
		val = v
	} else {
		return nil
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic("environment variable " + param + " is not a valid integer")
	}
	return &intVal
}

func (c *configuration) GetIntValue(param string) int {
	val := c.getIntValue(param)
	if val == nil {
		panic("environment variable " + param + " is not set")
	}
	return *val
}

func (c *configuration) GetStringValueDefault(param string, _default string) string {
	if v, ok := c.readSecretStoreVar(param); ok {
		return v
	}
	if val := os.Getenv(param); val != "" {
		return val
	}
	if v, ok := readDotEnvVar(param); ok {
		return v
	}
	return _default
}

func (c *configuration) GetIntValueDefault(param string, _default int) int {
	val := c.getIntValue(param)
	if val == nil {
		return _default
	}
	return *val
}

func (c *configuration) getFloat64Value(param string) *float64 {
	var val string
	if v, ok := c.readSecretStoreVar(param); ok {
		val = v
	} else if v := os.Getenv(param); v != "" {
		val = v
	} else if v, ok := readDotEnvVar(param); ok {
		val = v
	} else {
		return nil
	}
	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic("environment variable " + param + " is not a valid float64")
	}
	return &floatVal
}

func (c *configuration) GetFloat64Value(param string) float64 {
	val := c.getFloat64Value(param)
	if val == nil {
		panic("environment variable " + param + " is not set")
	}
	return *val
}

func (c *configuration) GetFloat64ValueDefault(param string, _default float64) float64 {
	val := c.getFloat64Value(param)
	if val == nil {
		return _default
	}
	return *val
}

func (c *configuration) IsDebugEnabled() bool {
	if v, ok := c.readSecretStoreVar("DEBUG"); ok {
		return strings.ToLower(strings.TrimSpace(v)) == "true"
	}
	if val := os.Getenv("DEBUG"); val != "" {
		return strings.ToLower(strings.TrimSpace(val)) == "true"
	}
	if v, ok := readDotEnvVar("DEBUG"); ok {
		return strings.ToLower(strings.TrimSpace(v)) == "true"
	}
	return false
}

func (c *configuration) IsProduction() bool {
	if v, ok := c.readSecretStoreVar("ENVIRONMENT"); ok {
		return strings.ToLower(strings.TrimSpace(v)) == "production"
	}
	if val := os.Getenv("ENVIRONMENT"); val != "" {
		return strings.ToLower(strings.TrimSpace(val)) == "production"
	}
	if v, ok := readDotEnvVar("ENVIRONMENT"); ok {
		return strings.ToLower(strings.TrimSpace(v)) == "production"
	}
	return false
}

func (c *configuration) readSecretStoreVar(param string) (string, bool) {
	repoName := c.repository
	context := c.context

	if repoName == "" || context == "" {
		return "", false
	}

	path := fmt.Sprintf("/mnt/secrets-store/_%s_%s_%s", repoName, context, param)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}

	val := strings.TrimSpace(string(data))
	return val, true
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
