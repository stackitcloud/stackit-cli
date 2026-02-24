package utils

import (
	"os"
	"strings"
)

func ParseScriptParams(params *map[string]string) (*map[string]string, error) { //nolint:gocritic // flag value is a map pointer
	if params == nil {
		return nil, nil
	}

	parsed := map[string]string{}
	for k, v := range *params {
		parsed[k] = v

		if k == "script" && strings.HasPrefix(v, "@{") && strings.HasSuffix(v, "}") {
			// Check if a script file path was specified, like: --params script=@{/tmp/test.sh}
			fileContents, err := os.ReadFile(v[2 : len(v)-1])
			if err != nil {
				return nil, err
			}
			parsed[k] = string(fileContents)
		}
	}

	return &parsed, nil
}
