package main

import "encoding/json"
import "fmt"
import "os"
import "os/exec"
import "github.com/mitchellh/go-homedir"
import "github.com/eugene-eeo/xfs/libxfs"

func handle(event libxfs.Event, d *libxfs.Dispatcher) error {
	switch event.Type {
	case libxfs.Create:
	case libxfs.Update:
		// dispatch event.Src | xfs-index set event.Src checksum
		checksum, err := libxfs.GetSHA256ChecksumFromFile(event.Src)
		if err != nil {
			return nil
		}
		mimetype, err := libxfs.MimetypeFromFile(event.Src)
		if err != nil {
			return nil
		}
		handler, found := d.Match(mimetype)
		if !found {
			return nil
		}
		file, err := os.Open(event.Src)
		if err != nil {
			return nil
		}
		defer file.Close()
		cmd := exec.Command(handler, event.Src)
		ind := exec.Command("bin/xfs-index", "set", event.Src, checksum)
		ind.Stdin, _ = cmd.StdoutPipe()
		ind.Start()
		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		err = ind.Wait()
		if err != nil {
			fmt.Println(err)
		}
	case libxfs.Delete:
		// xfs-index del event.Src
		exec.Command("bin/xfs-index", "del", event.Src).Run()
	case libxfs.Rename:
		// xfs-index move event.Src event.Dst
		exec.Command("bin/xfs-index", "move", event.Src, event.Dst).Run()
	}
	return nil
}

func main() {
	config, err := libxfs.NewConfig()
	if err != nil {
		panic(err)
	}
	dispatcher, err := libxfs.NewDispatcherFromJson(config.Dispatch)
	if err != nil {
		panic(err)
	}
	home, err := homedir.Expand("~/")
	if err != nil {
		panic(err)
	}
	d := json.NewDecoder(os.Stdin)
	for {
		e := libxfs.Event{}
		err = d.Decode(&e)
		if err == nil {
			fmt.Printf(
				"%s src=%s dst=%s\n",
				libxfs.PrettifyEventType(e.Type),
				libxfs.PrettifyPath(home, e.Src),
				libxfs.PrettifyPath(home, e.Dst),
			)
			go handle(e, dispatcher)
		} else {
			fmt.Println(err)
		}
	}
}
