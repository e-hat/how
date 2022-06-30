package main

import (
	"fmt"
	"io"
	"net/http"
)

func SayHi() {
	fmt.Println("hi!")
}

func repoEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello world!")
	io.WriteString(w, "Right back at ya!")
}

func main() {
	http.HandleFunc("/", repoEndpoint)

	err := http.ListenAndServe(":8000", nil)
	fmt.Println(err)
}
