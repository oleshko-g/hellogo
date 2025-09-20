package main

import (
	"errors"
	"flag"
	"os"
)

func parseArgs(args []string) (Config, error) {
	fs := flag.NewFlagSet("report", flag.ExitOnError)
	fs.SetOutput(os.Stderr)
	fs.String("input", "", "Path to the input file. Must be a non-empty string.")
	fs.Int("limit", 100, "Default: `100`. Must be between 1 and 1000 (inclusive)")
	verboseP := fs.Bool("verbose", false, "Default: `false`.")
	fs.String("format", "text", "Default: `text`. Allowed values: `text` or `json`.")

	err := fs.Parse(args)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err = cfg.setInput(fs); err != nil {
		return Config{}, err
	}
	if err = cfg.setLimit(fs); err != nil {
		return Config{}, err
	}

	cfg.Verbose = *verboseP

	if err = cfg.setFormat(fs); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

type Config struct {
	Input   string
	Limit   int
	Verbose bool
	Format  string
}

func (c *Config) setInput(fs *flag.FlagSet) error {
	input, ok := fs.Lookup("input").Value.(flag.Getter).Get().(string)
	if input == "" || !ok {
		return errors.New("missing -input")
	}

	c.Input = input
	return nil
}

func (c *Config) setLimit(fs *flag.FlagSet) error {
	value, ok := fs.Lookup("limit").Value.(flag.Getter)
	if !ok {
		return errors.New("invalid -limit: must be between 1 and 1000")
	}
	if limit, ok := value.Get().(int); ok {
		if !ok || (limit < 1 || limit > 1000) {
			return errors.New("invalid -limit: must be between 1 and 1000")
		}
		c.Limit = limit
	}
	return nil
}

func (c *Config) setFormat(fs *flag.FlagSet) error {
	if value, ok := fs.Lookup("format").Value.(flag.Getter); ok {
		if format, ok := value.Get().(string); ok {
			if format != "text" && format != "json" {
				return errors.New("invalid -format: must be 'text' or 'json'")
			}
			c.Format = format
		}
	}
	return nil
}
