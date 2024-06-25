package main

import (
	"fmt"
	"log"      // Importing the log package for logging
	"net/http" // Importing the net/http package for HTTP server functionality
)

// formHandler handles form submissions
func formHandler(w http.ResponseWriter, r *http.Request) {
	// ParseForm parses the raw query from the URL and updates r.Form
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	fmt.Fprintf(w, "")
	name := r.FormValue("text_input") // Retrieve the value of the 'name' parameter from the form
	// address := r.FormValue("address")         // Retrieve the value of the 'address' parameter from the form
	fmt.Fprintf(w, "%s\n", name) // Print the name to the response writer
	// fmt.Fprintf(w, "Address = %s\n", address) // Print the address to the response writer
}

// helloHandler handles requests to /hello endpoint
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusNotFound)
		return
	}
	// Respond with "hello!" for GET requests to /hello endpoint
	fmt.Fprintf(w, "hello!")
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	// Create a file server handler for serving static files from ./static directory
	http.Handle("/", fileServer)            // Handle requests to root path with fileServer handler
	http.HandleFunc("/form", formHandler)   // Handle requests to /form endpoint with formHandler function
	http.HandleFunc("/hello", helloHandler) // Handle requests to /hello endpoint with helloHandler function

	fmt.Printf("Starting server at port 8080\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		// Log any errors that occur during server startup and exit with error status
		log.Fatal(err)
	}
}
