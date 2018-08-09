package main

import "bufio"
import "encoding/json"
import "fmt"
import "os"
import "github.com/eugene-eeo/xfs/libxfs"

func main() {
	config, err := libxfs.NewConfig()
	if err != nil {
		panic(err)
	}
	dispatcher, err := libxfs.NewDispatcherFromJson(config.Dispatch)
	if err != nil {
		panic(err)
	}
	b := bufio.NewReader(os.Stdin)
	d := json.NewDecoder(b)
	for {
		e := &libxfs.Event{}
		err = d.Decode(e)
		if err == nil {
			fmt.Println("Event: ", e.Type, e.Src, e.Dst)
			if e.Type == libxfs.Create || e.Type == libxfs.Update {
				mimetype, err := libxfs.MimetypeFromFile(e.Src)
				if err != nil {
					panic(err)
				}
				fmt.Println(mimetype)
				fmt.Println(dispatcher.Match(mimetype))
			}
		}
	}
}
