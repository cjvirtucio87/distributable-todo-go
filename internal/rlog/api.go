package rlog

type Log interface {
	// Although entryId dictates the position in the
	// log, nextIndex dictates the starting point
	// for elements to be discarded in favor of
	// the parameter entries.
	AddEntries(nextIndex int, entries []*Entry) error
	Count() int
	Entry(idx int) *Entry
	Entries(start, end int) []*Entry
}

type Entry struct {
	// determines position in log
	Id      int
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
