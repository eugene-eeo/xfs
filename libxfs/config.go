package libxfs

import "io"
import "os"
import "encoding/json"
import "github.com/mitchellh/go-homedir"

type Config struct {
	Dispatch [][2]string `json:"dispatch"`
	Watch    []string    `json:"watch"`
}

func NewConfigFromReader(r io.Reader) (*Config, error) {
	config := &Config{}
	dec := json.NewDecoder(r)
	err := dec.Decode(config)
	if err != nil {
		return nil, err
	}
	if config.Watch == nil {
		config.Watch = []string{}
	}
	if config.Dispatch == nil {
		config.Dispatch = [][2]string{}
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
