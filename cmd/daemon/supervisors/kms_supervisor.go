package supervisors

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// KmsSupervisor - supervises internal KMS processes.
type KmsSupervisor struct {
	Supervisor

	// Reference to the daemon's context and wait group.
	//
	daemonCtx       context.Context
	daemonWaitGroup *sync.WaitGroup

	// Internal context and wait group for the KMS supervisor.
	//
	internalWaitGroup *sync.WaitGroup
	internalCtx       context.Context
	internalCancel    context.CancelFunc
}

// KmsSupervisorNew - constructor for KmsSupervisor.
func KmsSupervisorNew(daemonCtx context.Context, daemonWaitGroup *sync.WaitGroup) KmsSupervisor {

	internalCtx, internalCancel := context.WithCancel(context.Background())

	return KmsSupervisor{
		daemonCtx:         daemonCtx,
		daemonWaitGroup:   daemonWaitGroup,
		internalWaitGroup: &sync.WaitGroup{},
		internalCtx:       internalCtx,
		internalCancel:    internalCancel,
	}
}

// Start - starts the KMS API by creating its context and launching its internal process.
func (kA KmsSupervisor) Start() {
	defer kA.daemonWaitGroup.Done()

	// Enter KMS supervisor main loop.
	//
	kA.internalWaitGroup.Add(1)
	go kmsSupervisorMain(kA)

	// When the root context is done, stop the KMS supervisor.
	//
	<-kA.daemonCtx.Done()
	fmt.Println("KMS Supervisor stopping..")
	kA.Stop()
	kA.internalWaitGroup.Wait()
	fmt.Println("KMS Supervisor stopped")
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
			fmt.Println("KMS internal process stopped")
			return
		default:
			fmt.Println("KMS internal process is running")
			time.Sleep(1 * time.Second)
		}
	}
}
