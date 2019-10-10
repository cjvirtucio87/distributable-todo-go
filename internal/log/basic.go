package log

import "cjvirtucio87/distributed-todo-go/internal/dto"

type BasicLog struct {
	backend []dto.Entry
}

func (l *BasicLog) AddEntries(entries []dto.Entry) {
}

func (l *BasicLog) Count() int {
	return len(l.backend)
}

func (l *BasicLog) Entry(idx int) (dto.Entry, bool) {
	var e dto.Entry

	return e, false
}

func NewBasicLog() Log {
	return &BasicLog{
		backend: []dto.Entry{},
	}
}
