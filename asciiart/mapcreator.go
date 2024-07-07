package asciiart

import (
	"errors"
	"strings"
)

// MapCreator creates a map of ASCII art from a string
func MapCreator(s string) (map[rune][]string, error) {
	Map := make(map[rune][]string)
	var lines []string
	printableRune := rune(32) // ASCII space

	// Check if any art characters have been deleted from the bannerfile
	if len(s) != 6623 && len(s) != 5558 && len(s) != 7463 && len(s) != 6262 {
		return Map, errors.New("the bannerfile has been tampered with")
	}

	// Normalize line endings to \n for consistency
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines = strings.Split(s, "\n")

	// Number of lines per character in ASCII art
	numLinesPerChar := 8

	for i := 0; i < len(lines); i++ {
		// If the current line is empty and there are lines left to process
		if i+numLinesPerChar < len(lines) && lines[i] == "" {
			// Create a slice to store ASCII art lines for the current character
			artLines := []string{}
			// Iterate over 8 lines (assuming ASCII art is 8 lines tall)
			for j := 0; j < numLinesPerChar; j++ {
				// Append each line of ASCII art to the slice
				artLines = append(artLines, lines[i+1+j])
			}
			// Map the printable rune to its corresponding ASCII art
			Map[printableRune] = artLines
			// Increment the printable rune
			printableRune++
			// Skip the lines we've just processed
			i += numLinesPerChar
		}
	}

	return Map, nil
}
