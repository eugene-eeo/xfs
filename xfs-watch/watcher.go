package main

import "fmt"
import "github.com/rjeczalik/notify"
import "github.com/mitchellh/go-homedir"
import "github.com/eugene-eeo/xfs/libxfs"

func toEventType(event notify.Event) (t libxfs.EventType, ok bool) {
	ok = true
	switch event {
	case notify.Create:
		t = libxfs.Create
	case notify.Remove:
		t = libxfs.Delete
	case notify.Write:
		t = libxfs.Update
	default:
		ok = false
	}
	return
}

func main() {
	config, err := libxfs.NewConfig()
	if err != nil {
		panic(err)
	}
	agg := make(chan *libxfs.Event, 20)
	for _, path := range config.Watch {
		events := make(chan notify.EventInfo, 5)
		path, _ := homedir.Expand(path)
		if err := notify.Watch(path, events, notify.All); err != nil {
			return
		}
		defer notify.Stop(events)
		go (func() {
			var prev notify.EventInfo = nil
			for evt := range events {
				fmt.Println(evt.Event(), evt.Path())
				if evt.Event() == notify.Rename {
					if prev == nil {
						prev = evt
						continue
					}
					agg <- libxfs.NewEvent(libxfs.Rename, evt.Path(), prev.Path())
					prev = nil
					continue
				}
				ev_type, ok := toEventType(evt.Event())
				if !ok {
					continue
				}
				agg <- libxfs.NewEvent(ev_type, evt.Path(), "")
			}
		})()
	}
	for {
		event := <-agg
		fmt.Println(event.Type, event.Src, event.Dst)
	}
}
