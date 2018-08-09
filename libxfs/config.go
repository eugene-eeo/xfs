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
	Poll     int         `json:"poll"`
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
	if config.Watch == nil {
		config.Watch = []string{}
	}
	for i, x := range config.Watch {
		path, err := homedir.Expand(x)
		if err != nil {
			return nil, err
		}
		config.Watch[i] = path
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
