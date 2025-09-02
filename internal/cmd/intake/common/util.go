package common

import (
	"fmt"
	"strings"
)

// ParseLabels parses the labels flag value into a map.
// An empty string clears the labels, returning a pointer to an empty map.
func ParseLabels(labelsVal string) (map[string]string, error) {
	if labelsVal == "" {
		// User wants to clear labels
		return map[string]string{}, nil
	}

	// User provided labels, parse them
	parsedLabels := make(map[string]string)
	pairs := strings.Split(labelsVal, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 || kv[0] == "" {
			return nil, fmt.Errorf("invalid label format, expected key=value: %q", pair)
		}
		parsedLabels[kv[0]] = kv[1]
	}
	return parsedLabels, nil
}
