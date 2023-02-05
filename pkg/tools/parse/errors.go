package parse

func ErrorAsJSON(err error) map[string]interface{} {
	return map[string]interface{}{"error": err.Error()}
}
