package audit

type AuditorFile struct {
	Auditor
	storageDirectory string
	events           chan Event
}

func NewAuditorFile(storageDirectory string) *AuditorFile {
	return &AuditorFile{
		storageDirectory: storageDirectory,
		events:           make(chan Event, 100),
	}
}

// RecordEvent - records an auditing event.
func (aF *AuditorFile) RecordEvent(event string) error {
	// todo: implement logic
	return nil
}

// Persist - persists any buffered events to the storage.
func (aF *AuditorFile) Persist() error {
	// todo: implement logic
	return nil
}

// Close - closes the auditor and releases any resources.
func (aF *AuditorFile) Close() error {
	// todo: implement logic
	return nil
}
