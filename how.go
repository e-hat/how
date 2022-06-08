package main

import (
  "fmt"
  "io/ioutil"
  "os"
  "encoding/json"
)

type TopicEntry struct {
  Name string `json:"name"`
  Desc string `json:"desc"`
}

func repoPath() (string, bool) {
  if dirname, err := os.UserHomeDir(); err == nil {
    return fmt.Sprintf("%s/.how_repo/repo.json", dirname), true
  } else {
    return "", false
  }
}

func fetchRepo() (map[string]TopicEntry, bool) {
  path, ok := repoPath()
  if !ok {
    return make(map[string]TopicEntry), false
  }

  content, err := ioutil.ReadFile(path)
  if err != nil {
    return make(map[string]TopicEntry), false
  }

  var repo map[string]TopicEntry
  err = json.Unmarshal(content, &repo)
  if err != nil {
    return make(map[string]TopicEntry), false
  }

  return repo, true
}

func search(topic string) {
  repo, ok := fetchRepo()
  if !ok {
    fmt.Fprintln(os.Stderr, "error: could not load information repository")
  }

  // direct lookup
  if entry, ok := repo[topic]; ok {
    fmt.Printf("how information:\n\nName: %s\nDescription: %s\n", entry.Name, entry.Desc)
  } else {
    fmt.Println("needs fuzzy search")
  }
}

func writeRepo(repo *map[string]TopicEntry) bool {
  file, err := json.MarshalIndent(repo, "", " ")
  if err != nil {
    return false
  }

  path, ok := repoPath()
  if !ok {
    return false
  }

  err = ioutil.WriteFile(path, file, 0644)
  if err != nil {
    return false
  }

  return true
}

func write(entry TopicEntry) {
  repo, ok := fetchRepo()
  if !ok {
    fmt.Fprintln(os.Stderr, "error: could not load information repository")
  }

  repo[entry.Name] = entry

  if !writeRepo(&repo) {
    fmt.Fprintln(os.Stderr, "error: could not write repository")
  }
}

func usage() {
  fmt.Fprintln(os.Stderr, "usage: how <topic>|-- CMD\nwhere CMD is one of [help, write]")
}

type searchArgs struct {
  topic string
}

type writeArgs struct {
  name string
  desc string
}

type ArgType int64

const (
  ERROR ArgType = iota
  SEARCH
  ADD
)

func parseArgs() (ArgType, interface{}) {
  if len(os.Args) == 2 && os.Args[1] != "--" {
    return SEARCH, searchArgs { topic: os.Args[1] }
  } else if os.Args[1] == "--" {
    subargs := os.Args[2:]
    if len(subargs) == 0 {
      usage()
      return ERROR, nil
    }

    if subargs[0] == "write" && len(subargs) == 3 {
      return ADD, writeArgs { name: subargs[1], desc: subargs[2] }
    }
  } 

  usage()
  return ERROR, nil
}

func main() {
  type_, value := parseArgs()
  switch type_ {
  case ERROR:
    return
  case SEARCH: 
    args := value.(searchArgs)
    search(args.topic)
  case ADD:
    args := value.(writeArgs)
    write(TopicEntry{ Name: args.name, Desc: args.desc })
  }
}
