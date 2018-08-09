package libxfs

import "path/filepath"
import "github.com/blevesearch/bleve"

const BLEVE_FILE = "bleve"

type BleveEntry struct {
	Path     string
	Contents string
}

func GetBleveIndex(config *Config) (bleve.Index, error) {
	path := filepath.Join(config.DataDir, BLEVE_FILE)
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(path, mapping)
	if err == bleve.ErrorIndexPathExists {
		return bleve.Open(path)
	}
	return index, err
}

func (b *BleveEntry) Index(index bleve.Index) error {
	return index.Index(b.Path, b)
}

func GetBleveEntry(index bleve.Index, path string) (*BleveEntry, error) {
	doc, err := index.Document(path)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}
	contents := ""
	for _, field := range doc.Fields {
		if field.Name() == "Contents" {
			contents = string(field.Value())
			break
		}
	}
	return &BleveEntry{
		Path:     path,
		Contents: contents,
	}, nil
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func FilesMatchingQuery(index bleve.Index, query string, limit int) ([]string, error) {
	q := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	results, err := index.Search(q)
	if err != nil {
		return nil, err
	}
	r := make([]string, min(len(results.Hits), limit))
	for i, match := range results.Hits {
		if i == limit {
			break
		}
		r[i] = match.ID
	}
	return r, nil
}
