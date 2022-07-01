package main

import (
	"fmt"
	"os"

	"how/client"
	"how/repo"
	"how/server"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: how <topic>|-- CMD\nwhere CMD is one of [help, write, serve, push]")
}

type searchArgs struct {
	topic string
}

type writeArgs struct {
	name string
	desc string
}

type writeEditorArgs struct{}

type serveArgs struct{}

type pushArgs struct{}

type ArgType int64

const (
	ERROR ArgType = iota
	SEARCH
	WRITE
	WRITE_EDITOR
	SERVE_REPO
	PUSH
)

func parseArgs() (ArgType, interface{}) {
	if len(os.Args) == 2 && os.Args[1] != "--" {
		return SEARCH, searchArgs{topic: os.Args[1]}
	} else if len(os.Args) > 1 && os.Args[1] == "--" {
		subargs := os.Args[2:]
		if len(subargs) == 0 {
			usage()
			return ERROR, nil
		}

		if subargs[0] == "write" && len(subargs) == 3 {
			return WRITE, writeArgs{name: subargs[1], desc: subargs[2]}
		} else if subargs[0] == "write" && len(subargs) == 1 {
			return WRITE_EDITOR, writeEditorArgs{}
		} else if subargs[0] == "serve" && len(subargs) == 1 {
			return SERVE_REPO, serveArgs{}
		} else if subargs[0] == "push" && len(subargs) == 1 {
			return PUSH, pushArgs{}
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
		repo.Search(args.topic)
	case WRITE:
		args := value.(writeArgs)
		repo.WriteEntry(repo.TopicEntry{Name: args.name, Desc: args.desc})
	case WRITE_EDITOR:
		repo.WriteEditor()
	case SERVE_REPO:
		server.StartHowServer()
	case PUSH:
		client.Push()
	}
}
