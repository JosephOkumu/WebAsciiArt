package asciiart

import (
	"fmt"
	"strings"
)

// ValidateInput validates the user input.
func ValidateInput(input string) (string, error) {
	// Trim any leading or trailing whitespace
	input = strings.TrimSpace(input)

	// Check if the input is empty
	if input == "" {
		return "", fmt.Errorf("input cannot be empty")
	}

	// Check if the input contains only valid characters (printable ASCII)
	for _, char := range input {
		if char < 32 || char > 126 {
			return "", fmt.Errorf("invalid character in input: %q", char)
		}
	}

	return input, nil
}