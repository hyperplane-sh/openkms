package supervisors

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type CliAPISupervisor struct {
	Supervisor

	// Reference to the daemon's context and wait group.
	//
	daemonWaitGroup *sync.WaitGroup
	daemonCtx       context.Context

	// Internal context and wait group for the CLI API supervisor.
	//
	internalCtx       context.Context
	internalCancel    context.CancelFunc
	internalWaitGroup *sync.WaitGroup
}

// CliAPISupervisorNew - constructor for CliAPISupervisor.
func CliAPISupervisorNew(daemonCtx context.Context, daemonWaitGroup *sync.WaitGroup) CliAPISupervisor {
	internalCtx, internalCancel := context.WithCancel(context.Background())
	return CliAPISupervisor{
		daemonWaitGroup:   daemonWaitGroup,
		daemonCtx:         daemonCtx,
		internalCtx:       internalCtx,
		internalCancel:    internalCancel,
		internalWaitGroup: &sync.WaitGroup{},
	}
}

// Start - starts the CLI API.
func (cA CliAPISupervisor) Start() {
	defer cA.daemonWaitGroup.Done()

	cA.internalWaitGroup.Add(1)
	go cliAPISupervisorMain(cA)

	<-cA.daemonCtx.Done()
	fmt.Println("CLI Supervisor stopping..")
	cA.Stop()
	cA.internalWaitGroup.Wait()
	fmt.Println("CLI Supervisor stopped")
}

// Stop - stops the CLI API by cancelling its context.
func (cA CliAPISupervisor) Stop() {
	cA.internalCancel()
}

// Restart - restarts the CLI API by cancelling its current context and creating a new one.
func (cA CliAPISupervisor) Restart() {
	cA.internalCancel()
	time.Sleep(400 * time.Millisecond)

	cA.internalCtx, cA.internalCancel = context.WithCancel(context.Background())
	cA.daemonWaitGroup.Add(1)
	go cA.Start()
}

// cliAPISupervisorMain - main function for the CLI API daemon's internal process.
func cliAPISupervisorMain(cA CliAPISupervisor) {
	defer cA.internalWaitGroup.Done()
	for {
		select {
		case <-cA.internalCtx.Done():
			fmt.Println("CLI API internal process stopped")
			return
		default:
			fmt.Println("CLI API internal process is running")
			time.Sleep(1 * time.Second)
		}
	}
}
