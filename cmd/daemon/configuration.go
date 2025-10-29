package main

type DaemonConfiguration struct {
	CLI struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"CLI"`
}
