package log

import "cjvirtucio87/distributed-todo-go/internal/dto"

type Log interface {
	AddEntries(entryInfo dto.EntryInfo)
	Count() int
	Entry(idx int) (dto.Entry, bool)
	Entries(start, end int) ([]dto.Entry, bool)
}

func NewBasicLog(options ...func(*BasicLog)) Log {
	l := &BasicLog{}

	for _, o := range options {
		o(l)
	}

	return l
}

func WithBackend(backend []dto.Entry) func(*BasicLog) {
	return func(l *BasicLog) {
		l.backend = backend
	}
}
