package server

import (
	"io"
	"net/http"
  "fmt"
  "encoding/json"

  repoUtil "how/repo"
)

func repoEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
    repo, ok := repoUtil.Fetch()
    if !ok {
      fmt.Println("how: error: could not fetch repo")
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte("500 - Couldn't fetch repo from disk."))
      return
    }

    repoJson, _ := json.Marshal(repo)
    w.WriteHeader(http.StatusOK)
    w.Write(repoJson)
	case "PUT":
    body, err := io.ReadAll(r.Body)
    if err != nil {
      fmt.Println("how: error: could not read request")
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("400 - Couldn't read request"))
      return
    }

    repo, ok := repoUtil.Unmarshal(body)
    if !ok {
      fmt.Println("how: error: could not parse repo json")
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("400 - Couldn't parse the repo's json."))
      return
    }

    ok = repoUtil.Write(&repo)
    if !ok {
      fmt.Println("how: error: could not write repo")
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte("500 - Couldn't write the repo to disk."))
      return
    }

    w.WriteHeader(http.StatusOK)
	}
}

const PORT = 8000

func StartHowServer() {
  fmt.Println("how: starting server")
	http.HandleFunc("/", repoEndpoint)

  http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
