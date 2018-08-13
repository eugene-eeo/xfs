package main

import "os"
import "time"
import "strings"
import "encoding/json"
import "github.com/eugene-eeo/xfs/watcher"
import "github.com/eugene-eeo/xfs/libxfs"

func toEvent(evt watcher.Event) *libxfs.Event {
	ev_type := libxfs.EventType(-1)
	switch evt.Op {
	case watcher.Create:
		ev_type = libxfs.Create
	case watcher.Move:
		ev_type = libxfs.Rename
	case watcher.Rename:
		ev_type = libxfs.Rename
	case watcher.Remove:
		ev_type = libxfs.Delete
	case watcher.Write:
		ev_type = libxfs.Update
	}
	// Watcher emits events for each file when we do any of:
	//  $ rm -rf x
	//  $ mv x y
	//  $ touch x/a
	// So it is safe to ignore this event.
	if evt.IsDir() {
		return nil
	}
	return libxfs.NewEvent(ev_type, evt.Path, evt.Dst)
}

func addPaths(paths []string, w *watcher.Watcher) error {
	for _, path := range paths {
		if strings.HasSuffix(path, "/...") {
			path = strings.TrimSuffix(path, "/...")
			if err := w.AddRecursive(path); err != nil {
				return err
			}
			continue
		}
		if err := w.Add(path); err != nil {
			return err
		}
	}
	return nil
}

func addIgnores(paths []string, w *watcher.Watcher) error {
	for _, path := range paths {
		if err := w.Ignore(path); err != nil {
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
	w := watcher.New()
	w.FilterOps(
		watcher.Create,
		watcher.Write,
		watcher.Move,
		watcher.Rename,
		watcher.Remove,
		// TODO: handle chmod where we can't read the file any more
		// maybe map chmod => delete?
	)
	w.IgnoreHiddenFiles(true)
	err = addPaths(config.Watch, w)
	if err != nil {
		panic(err)
	}
	err = addIgnores(config.Ignore, w)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(os.Stdout)
	go func() {
		for {
			select {
			case evt := <-w.Event:
				ev := toEvent(evt)
				if ev != nil {
					enc.Encode(ev)
				}
			case err := <-w.Error:
				panic(err)
			case <-w.Closed:
				return
			}
		}
	}()
	w.Start(time.Second * time.Duration(config.Poll))
}
