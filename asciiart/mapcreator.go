package asciiart

import "strings"

// MapCreator creates a map of ASCII art from a string.
func MapCreator(s string) map[rune][]string {
	Map := make(map[rune][]string)
	printableRune := rune(32)
	lines := strings.Split(s, "\n")

	for i := 0; i < len(lines); i++ {
		// If the current line is empty and there are lines left to process
		if i+1 < len(lines) && lines[i] == "" {
			// Create a slice to store ASCII art lines for the current character
			artLines := []string{}
			// Iterate over 8 lines (assuming ASCII art is 8 lines tall)
			for j := 0; j < 8; j++ {
				// Append each line of ASCII art to the slice
				artLines = append(artLines, lines[i+1+j])
			}
			// Map the printable rune to its corresponding ASCII art
			Map[printableRune] = artLines
			// Increment the printable rune
			printableRune++
		}
	}
	// Map newline character to its corresponding ASCII art
	Map[rune(10)] = []string{"\n"}
	return Map
}
