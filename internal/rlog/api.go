package rlog

type Log interface {
	// Although entryId dictates the position in the
	// log, nextIndex dictates the starting point
	// for elements to be discarded in favor of
	// the parameter entries.
	AddEntries(nextIndex int, entries []*Entry) error
	// Retrieve a count of log entries.
	Count() int
	// Retrieve an entry at a specific index.
	Entry(idx int) *Entry
	// Retrieve a slice of entries from start to end.
	Entries(start, end int) []*Entry
}

type Entry struct {
	// The entry's position in the log.
	Id int
	// The command to be executed.
	Command string
}

func NewBasicLog(options ...func(*BasicLog)) Log {
	l := &BasicLog{}

	for _, o := range options {
		o(l)
	}

	return l
}

func WithBackend(backend []*Entry) func(*BasicLog) {
	return func(l *BasicLog) {
		l.backend = backend
	}
}
