package rlog

import (
	"cjvirtucio87/distributed-todo-go/internal/dto"
	"errors"
	"fmt"
)

type BasicLog struct {
	backend []dto.Entry
}

func (l *BasicLog) AddEntries(e dto.EntryInfo) error {
	l.backend = append(l.backend[:e.NextIndex], e.Entries...)

	return nil
}

func (l *BasicLog) Count() int {
	return len(l.backend)
}

func (l *BasicLog) Entry(idx int) (dto.Entry, error) {
	if e := l.backend[idx]; &e != nil {
		return e, nil
	} else {
		return e, errors.New(
			fmt.Sprintf(
				"could not retrieve entry for index %d",
				idx,
			),
		)
	}
}

func (l *BasicLog) Entries(start, end int) ([]dto.Entry, bool) {
	if start < 0 || end > l.Count() {
		return []dto.Entry{}, false
	}

	return l.backend[start:end], true
}
