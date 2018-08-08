package main

import "fmt"
import "path/filepath"
import "os"

import "github.com/eugene-eeo/xfs/libxfs"
import "github.com/docopt/docopt-go"
import bolt "github.com/coreos/bbolt"

func main() {
	usage := `
Usage:
	xfs-index set <path> <hash> [--dry-run]
	xfs-index get <path>
	xfs-index --help
	`
	arguments, _ := docopt.ParseDoc(usage)
	set := arguments["set"].(bool)
	get := arguments["get"].(bool)

	config, err := libxfs.NewConfig()
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll(config.DataDir, 0777); err != nil {
		panic(err)
	}

	db, err := bolt.Open(filepath.Join(config.DataDir, "bbolt"), 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = libxfs.InitDB(db)
	if err != nil {
		panic(err)
	}
	if set {
		path := libxfs.Path(arguments["<path>"].(string))
		hash := libxfs.Hash(arguments["<hash>"].(string))
		err = libxfs.SaveHash(db, path, hash)
		if err != nil {
			panic(err)
		}
	}
	if get {
		path := libxfs.Path(arguments["<path>"].(string))
		value, err := libxfs.GetHash(db, path)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(value))
	}
}
