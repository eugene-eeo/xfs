package libxfs

import "io"
import "os"
import "gopkg.in/h2non/filetype.v1"

func IsBinary(x []byte) bool {
	// Check if a buffer is binary or text
	// This does the same thing as git (check if the buffer
	// contains a NUL byte);
	// https://git.kernel.org/pub/scm/git/git.git/tree/xdiff-interface.c?id=HEAD#n202
	for _, b := range x {
		if b == 0 {
			return true
		}
	}
	return false
}

func Mimetype(buffer []byte) string {
	mimetype, _ := filetype.Match(buffer)
	if mimetype == filetype.Unknown {
		if !IsBinary(buffer) {
			return "application/x-text"
		}
		return "unknown"
	}
	return mimetype.MIME.Value
}

func MimetypeFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	b := make([]byte, 8000)
	n, err := file.Read(b)
	if err != nil && err != io.EOF {
		return "", err
	}
	return Mimetype(b[:n]), nil
}
