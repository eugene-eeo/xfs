package libxfs

import "io"
import "os"
import "encoding/json"
import "github.com/mitchellh/go-homedir"
import "github.com/DisposaBoy/JsonConfigReader"

type Config struct {
	DataDir  string      `json:"data_dir"`
	Dispatch [][2]string `json:"dispatch"`
	Watch    []string    `json:"watch"`
	Ignore   []string    `json:"ignore"`
	Poll     int         `json:"poll"`
}

func expandAll(paths []string) ([]string, error) {
	if paths == nil {
		return []string{}, nil
	}
	for i, x := range paths {
		path, err := homedir.Expand(x)
		if err != nil {
			return nil, err
		}
		paths[i] = path
	}
	return paths, nil
}

func NewConfigFromReader(r io.Reader) (*Config, error) {
	config := &Config{}
	dec := json.NewDecoder(JsonConfigReader.New(r))
	err := dec.Decode(config)
	if err != nil {
		return nil, err
	}
	if config.DataDir == "" {
		data_dir, err := homedir.Expand("~/.xfs")
		if err != nil {
			return nil, err
		}
		config.DataDir = data_dir
	}
	config.Watch, err = expandAll(config.Watch)
	if err != nil {
		return nil, err
	}
	config.Ignore, err = expandAll(config.Ignore)
	if err != nil {
		return nil, err
	}
	if config.Dispatch == nil {
		config.Dispatch = [][2]string{}
	}
	if config.Poll == 0 {
		config.Poll = 1
	}
	return config, nil
}

func NewConfig() (*Config, error) {
	path, err := homedir.Expand("~/.xfsrc")
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewConfigFromReader(file)
}
