package main

import (
  "fmt"
  "os"
)

func usage() {
  fmt.Fprintln(os.Stderr, "usage: how <topic>")
}

type searchArgs struct {
  term string
}

type ArgType int64

const (
  ERROR ArgType = iota
  SEARCH
)

func parseArgs() (ArgType, interface{}) {
  if len(os.Args) != 2 {
    usage()
    return ERROR, nil
  }

  return SEARCH, searchArgs { term: os.Args[1] }
}

func main() {
  type_, value := parseArgs()
  switch type_ {
  case ERROR:
    return
  case SEARCH: 
    args := value.(searchArgs)
    fmt.Printf("entered word: %s\n", args.term)
  }
}
