package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := "8080"
	fmt.Printf("✅ Server starting on http://localhost:%s\n", port)

	// A simple hello world handler to get started.
	// In a real app, you would wire up your handlers from the 'internal/adapters/http' package.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from %s!", "{{.ProjectName}}")
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}