package main

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	daemon = Daemon{
		waitGroup: sync.WaitGroup{},
	}
	openKMSConfigurationPath = getEnv("OPENKMS_CONFIG_PATH", "/etc/hyperplane/openkms/configs/openkms.yaml")
)

type Supervisor interface {
	// start - starts the supervised service.
	start()
	// stop - stops the supervised service.
	stop()
	// restart - restarts the supervised service.
	restart()
}

type Daemon struct {
	configuration DaemonConfiguration

	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc

	cliAPISupervisor cliAPISupervisor
	kmsSupervisor    kmsSupervisor
}

func init() {

	// First, make sure that the configuration file exists.
	//
	_, err := os.Stat(openKMSConfigurationPath)
	if os.IsNotExist(err) {
		slog.Error("Configuration file does not exist", "path", openKMSConfigurationPath)
		os.Exit(1)
	}

	// Load the configuration file.
	//
	content, err := os.ReadFile(openKMSConfigurationPath)
	if err != nil {
		slog.Error("Failed to read configuration file", "path", openKMSConfigurationPath, "error", err)
		os.Exit(1)
	}

	// Unmarshal the configuration file.
	//
	err = yaml.Unmarshal(content, &daemon.configuration)
	if err != nil {
		slog.Error("Failed to unmarshal configuration file", "path", openKMSConfigurationPath, "error", err)
		os.Exit(1)
	}

	daemon.ctx, daemon.cancel = context.WithCancel(context.Background())
}

// getEnv - retrieves the value of the environment variable named by the key.
func getEnv(key, def string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}
	return def
}

func main() {

	// Start supervisor for KMS API.
	//
	startSupervisor(daemon.kmsSupervisor)

	// Enable CLI API if enabled in configuration.
	//
	if daemon.configuration.CLI.Enabled == true {
		startSupervisor(daemon.cliAPISupervisor)
	}

	// Simulate running for 2 seconds before shutting down.
	//
	time.Sleep(2 * time.Second)
	daemon.cancel()

	// Wait for all goroutines to finish.
	//
	daemon.waitGroup.Wait()
}

func startSupervisor(s Supervisor) {
	daemon.waitGroup.Add(1)
	go s.start()
}
