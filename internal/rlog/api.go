package rlog

type Log interface {
	AddEntries(nextIndex int, entries []Entry) error
	Count() int
	Entry(idx int) *Entry
	Entries(start, end int) []Entry
}

type Entry struct {
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

func WithBackend(backend []Entry) func(*BasicLog) {
	return func(l *BasicLog) {
		l.backend = backend
	}
}
