package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"webart/asciiart"
)

// formHandler handles form submissions and generates ASCII art
func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() err: %v", err), http.StatusBadRequest)
		return
	}

	inputText := r.FormValue("text_input")
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

	asciiArtMap := asciiart.MapCreator(string(bannerContent))

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

	http.Redirect(w, r, "/display?art="+url.QueryEscape(asciiArtResult), http.StatusSeeOther)
}

// displayHandler handles displaying the generated ASCII art
func displayHandler(w http.ResponseWriter, r *http.Request) {
	asciiArtResult := r.URL.Query().Get("art")
	asciiArtResult = strings.ReplaceAll(asciiArtResult, "\n", "<br>")

	tmpl, err := template.ParseFiles("templates/display.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Art template.HTML
	}{
		Art: template.HTML(asciiArtResult),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "templates/index.html")
	})

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/display", displayHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
