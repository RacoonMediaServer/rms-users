package config

import "github.com/RacoonMediaServer/rms-packages/pkg/configuration"

type Security struct {
	Key          string
	Registration Registration
}

type Registration struct {
	Enabled bool
	Domains []string
}

// Configuration represents entire service configuration
type Configuration struct {
	Database configuration.Database
	Http     configuration.Http
	Monitor  configuration.Monitor
	Security Security
}

var config Configuration

// Load open and parses configuration file
func Load(configFilePath string) error {
	return configuration.Load(configFilePath, &config)
}

// Config returns loaded configuration
func Config() Configuration {
	return config
}
