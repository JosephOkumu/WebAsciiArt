package asciiart

import (
    "fmt"
    "regexp"
    "strings"
)

const maxInputLength = 1000 // Adjust this value as needed

// ArtRetriever returns the ASCII art corresponding to the input string using the provided map.
func ArtRetriever(s string, m map[rune][]string) (string, error) {
    if len(s) > maxInputLength {
        return "", fmt.Errorf("input too long: maximum allowed length is %d characters", maxInputLength)
    }

    var result strings.Builder
    re := regexp.MustCompile(`\\n`)
    s = re.ReplaceAllString(s, "\n")
    lines := strings.Split(s, "\n")

    for _, line := range lines {
        if line == "" {
            result.WriteString("\n")
        } else {
            for j := 0; j < 8; j++ {
                for _, char := range line {
                    if asciiArt, ok := m[char]; ok {
                        result.WriteString(asciiArt[j])
                    } else {
                        return "", fmt.Errorf("invalid input: %s", string(char))
                    }
                }
                result.WriteString("\n")
            }
        }
    }

    return strings.TrimSuffix(result.String(), "\n"), nil
}