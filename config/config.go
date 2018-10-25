package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mickep76/encoding"
	_ "github.com/mickep76/encoding/toml"
)

// Config struct.
type Config struct {
	LogLevel       string   `toml:"logLevel"`
	LogNoColor     bool     `toml:"logNoColor"`
	LogNoDate      bool     `toml:"logNoDate"`
	SkipInterfaces []string `toml:"skipInterfaces,omitempty"`
	SkipDrivers    []string `toml:"skipDrivers,omitempty"`
	Topic          string   `toml:"topic"`
	Brokers        []string `toml:"brokers"`
	Timeout        int      `toml:"timeout"`
	RateLimit      int      `toml:"rateLimit"`
	ComprAlgo      string   `toml:"comprAlgo"`
}

var codec encoding.Codec

func NewConfig() *Config {
	return &Config{
		LogLevel:       "debug",
		SkipInterfaces: []string{"lo"},
		SkipDrivers:    []string{"veth", "bridge"},
		Topic:          "ifwatch",
		Brokers:        []string{"kafka1", "kafka2", "kafka3"},
		Timeout:        3,
		RateLimit:      5,
		ComprAlgo:      "none",
	}
}

func (c *Config) Load(file string) error {
	// Check file type.
	if filepath.Ext(file) != ".toml" {
		return fmt.Errorf("unknown file format: %s", filepath.Ext(file))
	}

	// Read config file.
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// New encoding codec.
	if codec, err = encoding.NewCodec("toml"); err != nil {
		return err
	}

	// Decode config.
	if err := codec.Decode(b, c); err != nil {
		return err
	}

	return nil
}

func (c *Config) Print() {
	b, _ := codec.Encode(c)
	fmt.Println(string(b))
}
