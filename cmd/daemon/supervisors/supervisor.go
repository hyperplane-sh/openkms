package supervisors

import "github.com/hyperplane-sh/openkms/internal/audit"

type Supervisor interface {
	// Start - starts the supervised service.
	Start()
	// Stop - stops the supervised service.
	Stop()
	// Restart - restarts the supervised service.
	Restart()
	// Auditor  - returns the auditor used by the supervised service.
	Auditor() audit.Auditor
}
