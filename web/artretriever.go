package web

import (
	"fmt"
	"strings"
)

// ArtRetriever returns the ASCII art corresponding to the input string using the provided map.
func ArtRetriever(s string, m map[rune][]string) (string, error) {
	var result strings.Builder

	// Replace "\n" with actual newline characters
	// re := regexp.MustCompile(`\\n`)
	// s = re.ReplaceAllString(s, "\n")

	if EmptyOrNewlines(s) {
		return s, nil
	}

	lines := strings.Split(s, "\n")

	// Iterate over each line of the input string
	for ind := 0; ind < len(lines); ind++ {
		if lines[ind] == "" {
			// Add an empty line if the input line is empty
			result.WriteString("\n")
		} else {
			// Add ASCII art for each character in the input line
			for j := 0; j < 8; j++ {
				for _, char := range lines[ind] {
					if asciiArt, ok := m[char]; ok {
						// Add the corresponding ASCII art for the character
						result.WriteString(asciiArt[j])
					} else {
						return "", fmt.Errorf("error! invalid input: %s", string(char))
					}
				}
				result.WriteString("\n")
			}
		}
	}
	return result.String(), nil
}
