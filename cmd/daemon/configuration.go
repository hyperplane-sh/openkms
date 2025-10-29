package main

type DaemonConfiguration struct {
	CLI struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"CLI"`
	Auditing AuditingConfiguration `yaml:"Auditing"`
	KMS      KMSConfiguration      `yaml:"KMS"`
}

type AuditingConfiguration struct {
	Enabled bool   `yaml:"enabled"`
	Type    string `yaml:"type"`
	Storage struct {
		Directory string `yaml:"directory"`
	} `yaml:"storage"`
}
type KMSConfiguration struct{}
