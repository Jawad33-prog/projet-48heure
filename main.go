package main

import (
	"fmt"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	fmt.Println("Starting server at http://localhost:5555")
	if err := http.ListenAndServe(":5555", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
