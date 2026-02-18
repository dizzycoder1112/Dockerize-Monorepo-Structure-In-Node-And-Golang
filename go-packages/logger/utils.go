package logger

// ParseContext converts variadic context arguments to a map
// Supports two formats:
// 1. Key-value pairs: "key1", value1, "key2", value2
// 2. Map: map[string]any{"key1": value1, "key2": value2}
func ParseContext(context []any) map[string]any {
	if len(context) == 0 {
		return nil
	}

	// Check if first argument is a map
	if len(context) == 1 {
		if m, ok := context[0].(map[string]any); ok {
			return m
		}
	}

	// Parse as key-value pairs
	result := make(map[string]any)
	for i := 0; i < len(context); i += 2 {
		if i+1 >= len(context) {
			break
		}

		key, ok := context[i].(string)
		if !ok {
			continue
		}

		result[key] = context[i+1]
	}

	return result
}
