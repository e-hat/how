package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
  "io"

	repoUtil "how/repo"
	"how/server"
)

func Push() {
	envUrl := os.Getenv("HOW_URL")
	if len(envUrl) == 0 {
		fmt.Fprintln(os.Stderr, "error: the HOW_URL env var was empty")
    return
	}

	repo, ok := repoUtil.Fetch()
	if !ok {
		fmt.Fprintln(os.Stderr, "error: couldn't read repo from disk")
    return
	}

	repoJson, _ := json.Marshal(repo)

	url := fmt.Sprintf("http://%s:%d/", envUrl, server.PORT)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(repoJson))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    log.Fatal(err)
  }

  _, err = io.ReadAll(resp.Body)
  if err != nil {
    fmt.Fprintln(os.Stderr, "error: couldn't read response body")
    return
  }
}
