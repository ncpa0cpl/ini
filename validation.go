package ini

import "strings"

const DISALLOWED_KEY_CHARS = "?{}|&~![()^"

func isKeyValid(key string) bool {
	if key == "" {
		return false
	}

	if strings.ContainsAny(key, DISALLOWED_KEY_CHARS) {
		return false
	}

	return true
}
