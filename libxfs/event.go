package libxfs

type EventType int

const (
	Create = EventType(iota)
	Delete = EventType(iota)
	Update = EventType(iota)
	Rename = EventType(iota)
)

type Event struct {
	Type EventType
	Src  string
	Dst  string
}

func NewEvent(event_type EventType, src string, dst string) *Event {
	return &Event{
		Type: event_type,
		Src:  src,
		Dst:  dst,
	}
}
