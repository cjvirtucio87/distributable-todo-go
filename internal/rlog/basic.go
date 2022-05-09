package rlog

import "fmt"

type BasicLog struct {
	backend []*Entry
}

func (l *BasicLog) AddEntries(nextIndex int, entries []*Entry) error {
	backend := l.backend[:nextIndex]

	for _, entry := range entries {
		if entry.Id == nextIndex {
			backend = append(backend, entry)
		} else if entry.Id > nextIndex {
			return fmt.Errorf("entryId [%d] exceeds nextIndex [%d]", entry.Id, nextIndex)
		} else {
			backend[entry.Id] = entry
		}

		nextIndex++
	}

	l.backend = backend

	return nil
}

func (l *BasicLog) Count() int {
	return len(l.backend)
}

func (l *BasicLog) Entry(idx int) *Entry {
	return l.backend[idx]
}

func (l *BasicLog) Entries(start, end int) []*Entry {
	if start < 0 || end > l.Count() {
		return make([]*Entry, 0)
	}

	if end == -1 {
		return l.backend[start:]
	}

	return l.backend[start:end]
}
