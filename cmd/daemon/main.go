package main

import "syscall"

var (
	openKMSConfigurationPath = getEnv("OPENKMS_CONFIG_PATH", "/etc/hyperplane/openkms/configs/openkms.yaml")
)

func init() {

}

// getEnv - retrieves the value of the environment variable named by the key.
func getEnv(key, def string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}
	return def
}

func main() {
	println("OpenKMS Daemon")
}
