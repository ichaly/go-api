package util

func String(value string, defaultValue string) string {
	if len(value) > 0 {
		return value
	}
	return defaultValue
}
