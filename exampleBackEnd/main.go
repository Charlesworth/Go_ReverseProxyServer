package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hi there")
}

func main() {
	fmt.Println("Listening...")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":3000", nil)
}
