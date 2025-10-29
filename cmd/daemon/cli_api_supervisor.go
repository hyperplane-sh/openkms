package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type cliAPISupervisor struct {
	Supervisor
	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// start - starts the CLI API.
func (cA cliAPISupervisor) start() {
	defer daemon.waitGroup.Done()

	daemon.cliAPISupervisor.ctx, daemon.cliAPISupervisor.cancel = context.WithCancel(context.Background())
	daemon.cliAPISupervisor.waitGroup = sync.WaitGroup{}

	daemon.cliAPISupervisor.waitGroup.Add(1)
	go func() {
		defer daemon.cliAPISupervisor.waitGroup.Done()
		for {
			select {
			case <-daemon.cliAPISupervisor.ctx.Done():
				fmt.Println("CLI API internal process stopped")
				return
			default:
				fmt.Println("CLI API internal process is running")
				time.Sleep(1 * time.Second)
			}
		}
	}()

	for {
		select {
		case <-daemon.ctx.Done():
			fmt.Println("CLI Supervisor stopping..")
			daemon.cliAPISupervisor.stop()
			daemon.cliAPISupervisor.waitGroup.Wait()
			fmt.Println("CLI Supervisor stopped")
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

// stop - stops the CLI API by cancelling its context.
func (cA cliAPISupervisor) stop() {
	cA.cancel()
}

// restart - restarts the CLI API by cancelling its current context and creating a new one.
func (cA cliAPISupervisor) restart() {
	cA.cancel()
	time.Sleep(400 * time.Millisecond)

	cA.ctx, cA.cancel = context.WithCancel(context.Background())
	daemon.waitGroup.Add(1)
	go daemon.cliAPISupervisor.start()
}
