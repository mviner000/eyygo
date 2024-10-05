package utils

// BoolToString converts a boolean value to its string representation
func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
