package main

import "encoding/json"
import "os"
import "os/exec"
import "go.uber.org/zap"
import "github.com/mitchellh/go-homedir"
import "github.com/eugene-eeo/xfs/libxfs"

func logError(logger *zap.SugaredLogger, rid int, err error) {
	logger.Errorw(
		"error",
		"rid", rid,
		"err", err,
	)
}

func handle(event libxfs.Event, d *libxfs.Dispatcher) error {
	switch event.Type {
	case libxfs.Create:
	case libxfs.Update:
		// dispatch event.Src | xfs-index set event.Src checksum
		checksum, err := libxfs.GetSHA256ChecksumFromFile(event.Src)
		if err != nil {
			return err
		}
		mimetype, err := libxfs.MimetypeFromFile(event.Src)
		if err != nil {
			return err
		}
		handler, found := d.Match(mimetype)
		if !found {
			return nil
		}
		file, err := os.Open(event.Src)
		if err != nil {
			return err
		}
		defer file.Close()
		cmd := exec.Command(handler, event.Src)
		ind := exec.Command("bin/xfs-index", "set", event.Src, checksum)
		ind.Stdin, _ = cmd.StdoutPipe()
		_ = ind.Start()
		_ = cmd.Run()
		_ = ind.Wait()
	case libxfs.Delete:
		// xfs-index del event.Src
		err := exec.Command("bin/xfs-index", "del", event.Src).Run()
		if err != nil {
			return err
		}
	case libxfs.Rename:
		// xfs-index move event.Src event.Dst
		err := exec.Command("bin/xfs-index", "move", event.Src, event.Dst).Run()
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
	d := json.NewDecoder(os.Stdin)
	i := 0
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()
	for {
		e := libxfs.Event{}
		err = d.Decode(&e)
		if err == nil {
			i++
			rid := i
			go func() {
				sugar.Infow(
					"new request",
					"rid", rid,
					"type", libxfs.PrettifyEventType(e.Type),
					"src", libxfs.PrettifyPath(home, e.Src),
					"dst", libxfs.PrettifyPath(home, e.Dst),
				)
				err := handle(e, dispatcher)
				if err == nil {
					return
				}
				sugar.Errorw(
					"error handling request",
					"rid", rid,
					"err", err,
				)
			}()
		}
	}
}
