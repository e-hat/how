package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
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

func fuzzySearch(repo *map[string]TopicEntry, term string) []TopicEntry {
	keys := make([]string, 0, len(*repo))
	for k := range *repo {
		keys = append(keys, k)
	}

	ranks := []fuzzy.Rank(fuzzy.RankFind(term, keys))
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i].Distance > ranks[j].Distance
	})

	result := make([]TopicEntry, 0, len(ranks))
	for _, rank := range ranks {
		correspondingTopicEntry := (*repo)[rank.Target]
		result = append(result, correspondingTopicEntry)
	}

	return result
}

const SHORT_DESC_LEN = 20

func search(topic string) {
	repo, ok := fetchRepo()
	if !ok {
		fmt.Fprintln(os.Stderr, "error: could not load information repository")
	}

	// direct lookup
	if entry, ok := repo[topic]; ok {
		fmt.Printf("how information:\n\nName: %s\nDescription: %s\n", entry.Name, entry.Desc)
	} else {
		results := fuzzySearch(&repo, topic)
		if len(results) == 0 {
			fmt.Printf("0 results for search term '%s'\n", topic)
		} else {
			for i, entry := range results {
				shortDesc := entry.Desc[:int(math.Min(float64(len(entry.Desc)), float64(SHORT_DESC_LEN)))]
				fmt.Printf("...\n\nSearch Result #%d\nName: %s\nDescription: %s...\n\n...\n", i+1, entry.Name, shortDesc)
			}
		}
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
	WRITE
)

func parseArgs() (ArgType, interface{}) {
	if len(os.Args) == 2 && os.Args[1] != "--" {
		return SEARCH, searchArgs{topic: os.Args[1]}
	} else if os.Args[1] == "--" {
		subargs := os.Args[2:]
		if len(subargs) == 0 {
			usage()
			return ERROR, nil
		}

		if subargs[0] == "write" && len(subargs) == 3 {
			return WRITE, writeArgs{name: subargs[1], desc: subargs[2]}
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
	case WRITE:
		args := value.(writeArgs)
		write(TopicEntry{Name: args.name, Desc: args.desc})
	}
}
