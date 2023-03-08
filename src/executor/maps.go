package executor

func MapToString(e map[string]interface{}, v string) (string, bool) {
	if e == nil {
		return "", false
	}
	target, ok := e[v]
	if !ok {
		return "", false
	}
	val, ok := target.(string)
	if !ok {
		return "", false
	}
	return val, ok
}

func MapToFloat64(e map[string]interface{}, v string) (float64, bool) {
	if e == nil {
		return 0, false
	}
	target, ok := e[v]
	if !ok {
		return 0, false
	}
	val, ok := target.(float64)
	if !ok {
		return 0, false
	}
	return val, ok
}
