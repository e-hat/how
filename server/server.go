package server

import (
	"io"
	"log"
	"net/http"
	"os"
)

func repoEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		io.WriteString(w, "got your GET!")
	case "PUT":
		io.WriteString(w, "got your put")
	}
}

func StartHowServer() {
	logger := log.New(os.Stdout, "http:  ", log.LstdFlags)
	logger.Println("How Repository server starting...")

	http.HandleFunc("/", repoEndpoint)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
