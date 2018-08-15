package main

import "encoding/json"
import "os"
import "os/exec"
import "go.uber.org/zap"
import "github.com/mitchellh/go-homedir"
import "github.com/eugene-eeo/xfs/libxfs"

type Request struct {
	id int
	libxfs.Event
}

type HandlerError struct {
	stderr []byte
}

func (h *HandlerError) Error() string {
	return string(h.stderr)
}

func handle(event libxfs.Event, d *libxfs.Dispatcher) error {
	switch event.Type {
	case libxfs.Create:
	case libxfs.Update:
		// dispatch event.Src | xfs-index set event.Src
		mimetype, err := libxfs.MimetypeFromFile(event.Src)
		if err != nil {
			return err
		}
		handler, found := d.Match(mimetype)
		if !found {
			return nil
		}
		cmd := exec.Command(handler, event.Src)
		ind := exec.Command("xfs-index", "set", event.Src)
		ind.Stdin, err = cmd.StdoutPipe()
		if err != nil {
			return err
		}
		_ = ind.Start()
		_ = cmd.Run()
		_ = ind.Wait()
	case libxfs.Delete:
		// xfs-index del event.Src
		err := exec.Command("xfs-index", "del", event.Src).Run()
		if err != nil {
			return err
		}
	case libxfs.Rename:
		// xfs-index move event.Src event.Dst
		err := exec.Command("xfs-index", "move", event.Src, event.Dst).Run()
		if err != nil {
			return err
		}
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
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()
	requests := make(chan Request, 20)
	for n := 0; n < 4; n++ {
		go func() {
			for req := range requests {
				e := req.Event
				sugar.Infow("new request",
					"id", req.id,
					"type", libxfs.PrettifyEventType(e.Type),
					"src", libxfs.PrettifyPath(home, e.Src),
					"dst", libxfs.PrettifyPath(home, e.Dst))
				err := handle(e, dispatcher)
				if err != nil {
					sugar.Errorw("error handling request",
						"id", req.id,
						"err", err.Error())
				}
			}
		}()
	}
	d := json.NewDecoder(os.Stdin)
	i := 0
	for {
		e := libxfs.Event{}
		err = d.Decode(&e)
		if err == nil {
			i++
			requests <- Request{
				id:    i,
				Event: e,
			}
		}
	}
	close(requests)
}
