package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

type Supervisor interface {
	// start - starts the supervised service.
	start()
	// stop - stops the supervised service.
	stop()
	// restart - restarts the supervised service.
	restart()
}

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
	cliAPISupervisor cliAPISupervisor
	kmsSupervisor    kmsSupervisor
}

func handleSignalTermination() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	slog.Info("Termination signal received, shutting down...")
	daemon.cancel()
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

	// Load auditing if enabled.
	//
	if daemon.configuration.Auditing.Enabled == true {
		switch daemon.configuration.Auditing.Type {
		default:
			slog.Error("Unsupported auditing type", "type", daemon.configuration.Auditing.Type)
			os.Exit(1)
		case audit.TYPE_FILE:
			daemon.auditor = audit.NewAuditorFile(daemon.configuration.Auditing.Storage.Directory)
		}
	}

	if daemon.auditor != nil {
		daemon.waitGroup.Add(1)
		go func() {
			defer daemon.waitGroup.Done()

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					err := daemon.auditor.Persist()
					if err != nil {
						slog.Error("Failed to persist audit logs", "error", err)
					}
				case <-daemon.ctx.Done():
					// Persist any remaining logs before exiting.
					err := daemon.auditor.Persist()
					if err != nil {
						slog.Error("Failed to persist audit logs during shutdown", "error", err)
					}
					return
				}
			}
		}()
	}

	daemon.ctx, daemon.cancel = context.WithCancel(context.Background())

	go handleSignalTermination()
}

func main() {

	// Start supervisor for KMS API.
	//
	daemon.waitGroup.Add(1)
	go daemon.kmsSupervisor.start()

	// Enable CLI API if enabled in configuration.
	//
	if daemon.configuration.CLI.Enabled == true {
		daemon.waitGroup.Add(1)
		go daemon.cliAPISupervisor.start()
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
