package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"webart/asciiart"
)

// Add this struct to store the ASCII art result
type AsciiArtResult struct {
	Art string
}

// Create a map to store results temporarily
var resultStore = make(map[string]AsciiArtResult)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "templates/index.html")
}

func displayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Process form submission
		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("ParseForm() err: %v", err), http.StatusBadRequest)
			return
		}

		inputText := r.FormValue("text_input")
		if strings.Contains(inputText, "\r\n") {
			inputText = strings.ReplaceAll(inputText, "\r\n", "\n")
		}

		banner := r.FormValue("banner")
		var bannerFile string
		switch banner {
		case "standard":
			bannerFile = "banners/standard.txt"
		case "thinkertoy":
			bannerFile = "banners/thinkertoy.txt"
		case "shadow":
			bannerFile = "banners/shadow.txt"
		default:
			http.Error(w, "Invalid banner selection", http.StatusBadRequest)
			return
		}

		bannerContent, err := asciiart.ReadBannerFile(bannerFile)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading banner file: %v", err), http.StatusInternalServerError)
			return
		}

		asciiArtMap, err := asciiart.MapCreator(string(bannerContent))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error populating map: %v", err), http.StatusBadRequest)
			return
		}

		inputText, err = asciiart.ValidateInput(inputText)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error validating input: %v", err), http.StatusBadRequest)
			return
		}

		asciiArtResult, err := asciiart.ArtRetriever(inputText, asciiArtMap)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating ASCII art: %v", err), http.StatusInternalServerError)
			return
		}

		// Generate a unique ID for this result
		resultID := fmt.Sprintf("%d", len(resultStore) + 1)
		resultStore[resultID] = AsciiArtResult{Art: asciiArtResult}

		// Redirect to GET /display with the result ID
		http.Redirect(w, r, "/display?id=" + resultID, http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		// Display the result
		resultID := r.URL.Query().Get("id")
		result, exists := resultStore[resultID]
		if !exists {
			http.Error(w, "Result not found", http.StatusNotFound)
			return
		}

		tmpl, err := template.ParseFiles("templates/display.html")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
			return
		}

		data := struct {
			Art template.HTML
		}{
			Art: template.HTML(strings.ReplaceAll(result.Art, "\n", "<br>")),
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		}

		// Optional: remove the result from the store after displaying
		delete(resultStore, resultID)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/display", displayHandler)
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}