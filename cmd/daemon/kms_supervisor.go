package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type kmsSupervisor struct {
	Supervisor
	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// start - starts the KMS API by creating its context and launching its internal process.
func (kA kmsSupervisor) start() {
	defer daemon.waitGroup.Done()

	daemon.kmsSupervisor.ctx, daemon.kmsSupervisor.cancel = context.WithCancel(context.Background())
	daemon.kmsSupervisor.waitGroup = sync.WaitGroup{}

	daemon.kmsSupervisor.waitGroup.Add(1)
	go func() {
		defer daemon.kmsSupervisor.waitGroup.Done()
		for {
			select {
			case <-daemon.kmsSupervisor.ctx.Done():
				fmt.Println("KMS API stopped")
				return
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	for {
		select {
		case <-daemon.ctx.Done():
			fmt.Println("KMS Supervisor stopping..")
			daemon.kmsSupervisor.stop()
			daemon.kmsSupervisor.waitGroup.Wait()
			fmt.Println("KMS Supervisor stopped")
			return
		}
	}
}

// stop - stops the KMS API by cancelling its context.
func (kA kmsSupervisor) stop() {
	kA.cancel()
}

// restart - restarts the KMS API by cancelling its current context and creating a new one.
func (kA kmsSupervisor) restart() {
	kA.cancel()
	time.Sleep(400 * time.Millisecond)

	kA.ctx, kA.cancel = context.WithCancel(context.Background())

	daemon.waitGroup.Add(1)
	go daemon.kmsSupervisor.start()
}
