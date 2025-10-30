package supervisors

type Supervisor interface {
	// Start - starts the supervised service.
	Start()
	// Stop - stops the supervised service.
	Stop()
	// Restart - restarts the supervised service.
	Restart()
}
