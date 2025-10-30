package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hyperplane-sh/openkms/cmd/daemon/supervisors"
	"github.com/hyperplane-sh/openkms/internal/audit"
	"gopkg.in/yaml.v3"
)

var (
	daemon = Daemon{
		configurationLock: sync.RWMutex{},
		waitGroup:         sync.WaitGroup{},
	}
	openKMSConfigurationPath = getEnv("OPENKMS_CONFIG_PATH", "/etc/hyperplane/openkms/configs/openkms.yaml")
)

type Daemon struct {
	configuration     DaemonConfiguration
	configurationLock sync.RWMutex // protects access to the configuration when loading or reloading.

	// Auditing related fields.
	//
	auditor audit.Auditor

	// Root context and wait group for the daemon.
	//
	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc

	// Supervisors
	//
	cliAPISupervisor supervisors.CliAPISupervisor
	kmsSupervisor    supervisors.KmsSupervisor
}

func handleSignalTermination() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	slog.Info("Termination signal received, shutting down...")
	daemon.cancel()

	daemon.waitGroup.Wait()

	slog.Info("Shutdown complete")
	os.Exit(0)
}

func init() {

	// First, make sure that the configuration file exists.
	//
	daemon.configurationLock.Lock() // Lock the configuration while loading.

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
	daemon.configurationLock.Unlock()

	daemon.ctx, daemon.cancel = context.WithCancel(context.Background())

	// Load auditing if enabled.
	//
	if daemon.configuration.Auditing.Enabled == true {
		switch daemon.configuration.Auditing.Type {
		default:
			slog.Error("Unsupported auditing type", "type", daemon.configuration.Auditing.Type)
			os.Exit(1)
		case audit.TYPE_FILE:
			daemon.auditor = audit.NewAuditorFile(daemon.configuration.Auditing.Storage.Directory, daemon.ctx)
		}
	}

	if daemon.auditor != nil {
		daemon.waitGroup.Add(1)
		go func() {
			defer daemon.waitGroup.Done()
			err = daemon.auditor.Persist()
			if err != nil {
				slog.Error("Failed to persist auditing events", "error", err)
				os.Exit(1)
			}
		}()
	}

	go handleSignalTermination()
}

func main() {

	daemon.auditor.RecordEvent(audit.NewEvent(
		audit.LEVEL_INFO,
		"DAEMON",
		audit.TOPIC_LIFECYCLE,
		"OpenKMS daemon is starting up",
		map[string]string{},
	))

	// Start supervisor for KMS API.
	//
	daemon.waitGroup.Add(1)
	daemon.kmsSupervisor = supervisors.KmsSupervisorNew(daemon.ctx, &daemon.waitGroup, daemon.auditor)
	go daemon.kmsSupervisor.Start()

	// Enable CLI API if enabled in configuration.
	//
	if daemon.configuration.CLI.Enabled == true {
		daemon.waitGroup.Add(1)
		daemon.cliAPISupervisor = supervisors.CliAPISupervisorNew(daemon.ctx, &daemon.waitGroup, daemon.auditor)
		go daemon.cliAPISupervisor.Start()
	}

	// Wait for all goroutines to finish.
	//
	daemon.waitGroup.Wait()
}

// getEnv - retrieves the value of the environment variable named by the key.
func getEnv(key, def string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}
	return def
}
