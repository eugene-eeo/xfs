package libxfs

import "strings"

func PrettifyPath(home, path string) string {
	if strings.HasPrefix(path, home) {
		return strings.Replace(path, home, "~", 1)
	}
	return path
}

func PrettifyEventType(t EventType) string {
	switch t {
	case Rename:
		return "Rename"
	case Update:
		return "Update"
	case Delete:
		return "Delete"
	case Create:
		return "Create"
	default:
		return "Unknown"
	}
}
