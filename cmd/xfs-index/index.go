package main

import "io/ioutil"
import "fmt"
import "path/filepath"
import "os"

import "github.com/eugene-eeo/xfs/libxfs"
import "github.com/docopt/docopt-go"
import "github.com/blevesearch/bleve"
import bolt "github.com/coreos/bbolt"

func indexStdin(path libxfs.Path, index bleve.Index) error {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	entry := libxfs.BleveEntry{
		Path:     string(path),
		Contents: string(b),
	}
	return entry.Index(index)
}

func main() {
	usage := `
Usage:
	xfs-index set <path> <hash> [--dry-run]
	xfs-index get <path>
	xfs-index del <path>
	xfs-index move <src> <dst>
	xfs-index --help
	`
	arguments, _ := docopt.ParseDoc(usage)
	set := arguments["set"].(bool)
	get := arguments["get"].(bool)
	del := arguments["del"].(bool)
	move := arguments["move"].(bool)

	config, err := libxfs.NewConfig()
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll(config.DataDir, 0777); err != nil {
		panic(err)
	}

	index, err := libxfs.GetBleveIndex(config)
	if err != nil {
		panic(err)
	}

	db, err := bolt.Open(filepath.Join(config.DataDir, libxfs.BBOLT_FILENAME), 0600, nil)
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
		err = indexStdin(path, index)
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
		fmt.Println()
		entry, err := libxfs.GetBleveEntry(index, string(path))
		if err != nil {
			panic(err)
		}
		if entry != nil {
			fmt.Println(entry.Contents)
		}
	}
	if del {
		path := libxfs.Path(arguments["<path>"].(string))
		err := libxfs.DelHash(db, path)
		if err != nil {
			panic(err)
		}
		if err := libxfs.DelBleveEntry(index, string(path)); err != nil {
			panic(err)
		}
	}
	if move {
		src := libxfs.Path(arguments["<src>"].(string))
		dst := libxfs.Path(arguments["<dst>"].(string))
		err := libxfs.MoveHash(db, src, dst)
		if err != nil {
			panic(err)
		}
		if err := libxfs.MoveBleveEntry(index, string(src), string(dst)); err != nil {
			panic(err)
		}
	}
}
