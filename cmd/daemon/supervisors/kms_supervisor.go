package supervisors

import (
	"context"
	"sync"
	"time"

	"github.com/hyperplane-sh/openkms/internal/audit"
)

const (
	KMS_SUPERVISOR_AUDIT_GROUP = "KMS-SUPERVISOR"
)

// KmsSupervisor - supervises internal KMS processes.
type KmsSupervisor struct {
	Supervisor

	// Reference to the daemon's context and wait group.
	//
	daemonWaitGroup *sync.WaitGroup
	daemonCtx       context.Context

	// Auditor
	//
	auditor audit.Auditor

	// Internal context and wait group for the KMS supervisor.
	//
	internalWaitGroup *sync.WaitGroup
	internalCtx       context.Context
	internalCancel    context.CancelFunc
}

// KmsSupervisorNew - constructor for KmsSupervisor.
func KmsSupervisorNew(daemonCtx context.Context, daemonWaitGroup *sync.WaitGroup, auditor audit.Auditor) KmsSupervisor {

	internalCtx, internalCancel := context.WithCancel(context.Background())

	return KmsSupervisor{
		daemonWaitGroup:   daemonWaitGroup,
		daemonCtx:         daemonCtx,
		auditor:           auditor,
		internalWaitGroup: &sync.WaitGroup{},
		internalCtx:       internalCtx,
		internalCancel:    internalCancel,
	}
}

// Start - starts the KMS API by creating its context and launching its internal process.
func (kA KmsSupervisor) Start() {
	defer kA.daemonWaitGroup.Done()

	kA.auditor.RecordEvent(audit.NewEvent(
		audit.LEVEL_INFO,
		KMS_SUPERVISOR_AUDIT_GROUP,
		audit.TOPIC_LIFECYCLE,
		"KMS Supervisor starting",
		map[string]string{},
	))

	// Enter KMS supervisor main loop.
	//
	kA.internalWaitGroup.Add(1)
	go kmsSupervisorMain(kA)

	// When the root context is done, stop the KMS supervisor.
	//
	<-kA.daemonCtx.Done()
	kA.auditor.RecordEvent(audit.NewEvent(
		audit.LEVEL_INFO,
		KMS_SUPERVISOR_AUDIT_GROUP,
		audit.TOPIC_LIFECYCLE,
		"KMS Supervisor stopping",
		map[string]string{},
	))

	kA.Stop()
	kA.internalWaitGroup.Wait()
}

// Stop - stops the KMS API by cancelling its context.
func (kA KmsSupervisor) Stop() {
	kA.internalCancel()
}

// Restart - restarts the KMS API by cancelling its current context and creating a new one.
func (kA KmsSupervisor) Restart() {
	kA.internalCancel()
	time.Sleep(400 * time.Millisecond)

	kA.internalCtx, kA.internalCancel = context.WithCancel(context.Background())

	kA.daemonWaitGroup.Add(1)
	go kA.Start()
}

// kmsSupervisorMain - main function for the KMS daemon's internal process.
func kmsSupervisorMain(kA KmsSupervisor) {
	defer kA.internalWaitGroup.Done()
	for {
		select {
		case <-kA.internalCtx.Done():
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
