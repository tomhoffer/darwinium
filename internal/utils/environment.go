package utils

import (
	"os"
	"strings"
)

// IsTestEnvironment checks if the code is running in a test environment
func IsTestEnvironment() bool {
	// Check for test flags or environment variables
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}
