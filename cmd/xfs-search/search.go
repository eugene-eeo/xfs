package main

import "fmt"
import "strconv"
import "github.com/docopt/docopt-go"
import "github.com/mitchellh/go-homedir"
import "github.com/eugene-eeo/xfs/libxfs"

func main() {
	arguments, _ := docopt.ParseDoc(`
Usage:
	xfs-search <query> [--limit=<n>] [--pretty]
	xfs-search --help
	`)
	pretty := arguments["--pretty"].(bool)
	query := arguments["<query>"].(string)
	limit := 0
	limitStr := arguments["--limit"]
	if limitStr == nil {
		limit = 20
	} else {
		l, err := strconv.Atoi(limitStr.(string))
		if err != nil {
			panic("Expected --limit to be number")
		}
		limit = l
	}
	home, err := homedir.Expand("~/")
	if err != nil {
		panic(err)
	}
	config, err := libxfs.NewConfig()
	if err != nil {
		panic(err)
	}
	index, err := libxfs.GetBleveIndex(config)
	if err != nil {
		panic(err)
	}
	results, err := libxfs.FilesMatchingQuery(index, query, limit)
	if err != nil {
		panic(err)
	}
	if pretty {
		for _, f := range results {
			fmt.Println(libxfs.PrettifyPath(home, f))
		}
	} else {
		for _, f := range results {
			fmt.Println(f)
		}
	}
}
