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

	return input, nil
}