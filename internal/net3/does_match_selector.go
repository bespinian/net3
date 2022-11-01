package net3

// doesMatchSelector checks if a set of labels matches a selector.
func doesMatchSelector(selector, labels map[string]string) bool {
	for k, v := range selector {
		if labels[k] != v {
			return false
		}
	}
	return true
}
