package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

func Pull() {
	envUrl := os.Getenv("HOW_URL")
	if len(envUrl) == 0 {
		fmt.Fprintln(os.Stderr, "error: the HOW_URL env var was empty")
		return
	}

	url := fmt.Sprintf("http://%s:%d/", envUrl, server.PORT)
	res, err := http.Get(url)
  if err != nil {
    fmt.Println(err)
    return
  }

  repoJson, err := io.ReadAll(res.Body)
  if err != nil {
    fmt.Println("error: couldn't read response body")
    return
  }

  repo, ok := repoUtil.Unmarshal(repoJson)
  if !ok {
    fmt.Println("error: did not receive valid repo as response")
    return
  }

  ok = repoUtil.Write(&repo)
  if !ok {
    fmt.Println("error: failed to write repo to disk")
  }
}
