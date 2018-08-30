package main

import "io/ioutil"
import "fmt"
import "os"

import "github.com/eugene-eeo/xfs/libxfs"
import "github.com/docopt/docopt-go"
import "github.com/blevesearch/bleve"

func indexStdin(path string, index bleve.Index) error {
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
	xfs-index set <path> [--dry-run]
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

	if set {
		path := arguments["<path>"].(string)
		err = indexStdin(path, index)
		if err != nil {
			panic(err)
		}
	}
	if get {
		path := arguments["<path>"].(string)
		entry, err := libxfs.GetBleveEntry(index, path)
		if err != nil {
			panic(err)
		}
		if entry != nil {
			fmt.Println(path, "OK")
			fmt.Println("--------")
			fmt.Println(entry.Contents)
		}
	}
	if del {
		path := arguments["<path>"].(string)
		if err := libxfs.DelBleveEntry(index, path); err != nil {
			panic(err)
		}
	}
	if move {
		src := arguments["<src>"].(string)
		dst := arguments["<dst>"].(string)
		if err := libxfs.MoveBleveEntry(index, src, dst); err != nil {
			panic(err)
		}
	}
}
