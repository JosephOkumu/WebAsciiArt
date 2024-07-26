package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"asciiart/web"
)

type Ascii struct {
	Result string
}

func main() {
	if len(os.Args) != 1 {
		return
	}

	// Create a new HTTP request multiplexer
	mux := http.NewServeMux()

	// Handle static file requests
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Set up route handlers
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/ascii-art", asciiArtHandler)

	fmt.Println("Starting the Server at port 8080")

	// Create an HTTP server with the specified address and handler
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the server and listen for incoming requests
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}

// homeHandler handles requests to the root path ("/")
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Print the received request path
	fmt.Printf("Received request to %s\n", r.URL.Path)

	// Check if the request path is the root path
	if r.URL.Path != "/" {
		renderError(w, http.StatusNotFound, "", "Not Found", "The page you are looking for could not be found.")
		return
	}

	// Handle GET requests to the root path
	if r.Method == http.MethodGet {
		// Parse the index.html template file
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			renderError(w, http.StatusNotFound, err.Error(), "Not Found", "The template you are looking for could not be found.")
			return
		}

		// Execute the template and write the response
		t.Execute(w, nil)
	} else {
		// Handle other request methods with a bad request error
		renderError(w, http.StatusBadRequest, "", "Bad Request", "Your request could not be processed.")
	}
}

// asciiArtHandler handles requests to the "/ascii-art" path
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	// Print the received request path
	fmt.Printf("Received request to %s\n", r.URL.Path)

	// Check if the request path is the /ascii-art path
	if r.URL.Path != "/ascii-art" {
		renderError(w, http.StatusNotFound, "", "Not Found", "The page you are looking for could not be found.")
		return
	}

	// Handle POST requests to the /ascii-art path
	if r.Method == http.MethodPost {
		// Get the text and banner values from the form
		text := r.FormValue("text")
		if strings.Contains(text, "\r\n") {
			text = strings.ReplaceAll(text, "\r\n", "\n")
		}
		banner := r.FormValue("banner")

		// Generate ASCII art using the provided text and banner
		result, err := generateAsciiArt(text, banner)
		if err != nil {
			if strings.Contains(err.Error(), "invalid input") {
				renderError(w, http.StatusInternalServerError, err.Error(), "Internal Server Error", "The input text contains non-ASCII characters that could not be converted to ASCII art")
				return
			}
			renderError(w, http.StatusNotFound, err.Error(), "Not Found", "The banner you are looking for could not be found.")
			return
		}

		// Parse the result.html template file
		t, err := template.ParseFiles("templates/result.html")
		if err != nil {
			renderError(w, http.StatusNotFound, err.Error(), "Not Found", "The template you are looking for could not be found.")
			return
		}

		// Create an Ascii struct with the generated ASCII art
		ascii := Ascii{
			Result: result,
		}

		// Execute the template and write the response
		t.Execute(w, ascii)
	} else {
		// Handle other request methods with a bad request error
		renderError(w, http.StatusBadRequest, "", "Bad Request", "Your request could not be processed.")
	}
}

// generateAsciiArt generates ASCII art for the given text and banner
func generateAsciiArt(text, banner string) (string, error) {
	// Read the banner file corresponding to the specified banner
	bannerFile, err := web.ReadBannerFile("./banners/" + banner + ".txt")
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	// Create a map of each printable ASCII rune to its ASCII art representation
	runeAsciiArtMap := web.MapCreator(bannerFile)

	// Generate ASCII art for the input text using the created map
	artText, err := web.ArtRetriever(text, runeAsciiArtMap)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return artText, nil
}

// renderError renders an error response with the specified status code and message
func renderError(w http.ResponseWriter, status int, errorMessage string, message string, information string) {
	// Set the response status code
	w.WriteHeader(status)

	var templateFile string = "templates/status_codes.html"

	// Parse the error template file
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a data struct with the error message
	data := struct {
		Status  int
		Message string
		Error   string
		Info    string
	}{
		Status:  status,
		Message: message,
		Error:   errorMessage,
		Info:    information,
	}

	// Execute the template and write the response
	t.Execute(w, data)
}
