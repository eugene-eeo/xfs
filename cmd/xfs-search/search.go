package main

import "fmt"
import "strconv"
import "github.com/docopt/docopt-go"
import "github.com/eugene-eeo/xfs/libxfs"

func main() {
	arguments, _ := docopt.ParseDoc(`
Usage:
	xfs-search <query> [--limit=<n>]
	xfs-search --help
	`)
	query := arguments["<query>"].(string)
	limitStr := arguments["--limit"].(string)
	limit := 0
	if limitStr == "" {
		limit = 20
	} else {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			panic("Expected --limit to be number")
		}
		limit = l
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
	for _, f := range results {
		fmt.Println(f)
	}
}
