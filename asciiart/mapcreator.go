package asciiart

import (
	"errors"
	"strings"
)

// MapCreator creates a map of ASCII art from a string
func MapCreator(s string) (map[rune][]string, error) {
	Map := make(map[rune][]string)
	var lines []string
	printableRune := rune(32)

	// Check if any art characters have been deleted from the bannerfile
	if len(s) != 6623 && len(s) != 5558 && len(s) != 7463 && len(s) != 6262 {
		return Map, errors.New("the bannerfile has been tampered with")
	}

	// Check for how lines are split in the banner file
	if strings.ContainsRune(s, '\r') {
		lines = strings.Split(s, "\r\n")
	} else {
		lines = strings.Split(s, "\n")
	}

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
	return Map, nil
}