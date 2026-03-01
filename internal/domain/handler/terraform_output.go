package handler

import "strings"

// These error message indicate a transient network issue that may be solved by retrying.
var FlakyNetworkStrings []string = []string{"network is unreachable", "handshake timeout"}

func IsFlakyNetworkError(output string) bool {
	for _, flakyString := range FlakyNetworkStrings {
		if strings.Contains(output, flakyString) {
			return true
		}
	}

	return false
}
