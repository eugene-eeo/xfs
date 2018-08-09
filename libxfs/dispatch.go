package libxfs

import "github.com/gobwas/glob"

type Entry struct {
	glob    glob.Glob
	handler string
}

type Dispatcher struct {
	entries []Entry
}

func NewDispatcherFromJson(data [][2]string) (*Dispatcher, error) {
	entries := make([]Entry, len(data))
	for i, entry := range data {
		pattern := entry[0]
		handler := entry[1]
		g, err := glob.Compile(pattern)
		if err != nil {
			return nil, err
		}
		entries[i] = Entry{g, handler}
	}
	return &Dispatcher{entries}, nil
}

func (d *Dispatcher) Match(mimetype string) (handler string, found bool) {
	for _, entry := range d.entries {
		if entry.glob.Match(mimetype) {
			return entry.handler, true
		}
	}
	return "", false
}
