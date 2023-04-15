package utils

func InStringSlice(element string, slice []string) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}

	return false
}
