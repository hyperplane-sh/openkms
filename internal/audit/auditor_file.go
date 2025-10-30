package audit

import (
	"context"
	"os"
	"sync"
	"time"
)

type AuditorFile struct {
	Auditor
	daemonCtx        context.Context
	storageDirectory string
	events           chan Event
	eventsWriteLock  sync.Mutex
}

func NewAuditorFile(storageDirectory string, daemonCtx context.Context) *AuditorFile {
	return &AuditorFile{
		daemonCtx:        daemonCtx,
		storageDirectory: storageDirectory,
		events:           make(chan Event, 100),
		eventsWriteLock:  sync.Mutex{},
	}
}

// RecordEvent - records an auditing event.
func (aF *AuditorFile) RecordEvent(event Event) error {
	aF.eventsWriteLock.Lock()
	defer aF.eventsWriteLock.Unlock()

	// When the channel is full, wait and try again.
	//
	if len(aF.events) == cap(aF.events) {
		time.Sleep(100 * time.Millisecond)
		return aF.RecordEvent(event)
	}

	aF.events <- event

	return nil
}

// Persist - persists any buffered events to the storage.
func (aF *AuditorFile) Persist() error {

	logFile := aF.getFileName()
	logPath := aF.storageDirectory + "/" + logFile

	// Create log file (if not exists) and open for appending.
	//
	_, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		_, err := os.Create(logPath)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for {
		select {
		case event := <-aF.events:
			_, err := f.WriteString(event.ToString() + "\n")
			if err != nil {
				return err
			}
		case <-aF.daemonCtx.Done():
			// Drain remaining events before exiting.
			//
			for len(aF.events) > 0 {
				event := <-aF.events
				_, err := f.WriteString(event.ToString() + "\n")
				if err != nil {
					return err
				}
			}
			return nil
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil
}

// Close - closes the auditor and releases any resources.
func (aF *AuditorFile) Close() error {
	// todo: implement logic
	return nil
}

// getFileName - generates a file name based on the current date.
func (aF *AuditorFile) getFileName() string {
	return time.Now().Format(time.RFC3339) + "_audit.log"
}
