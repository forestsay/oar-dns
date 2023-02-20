package util

import "os"

func GetEnvVar(varName, defaultVal string) string {
	ret, exist := os.LookupEnv(varName)
	if !exist || ret == "" {
		return defaultVal
	}
	return ret
}
