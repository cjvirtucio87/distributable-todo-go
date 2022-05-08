package rlog

type BasicLog struct {
	backend []Entry
}

func (l *BasicLog) AddEntries(nextIndex int, entries []Entry) error {
	l.backend = append(l.backend[:nextIndex], entries...)

	return nil
}

func (l *BasicLog) Count() int {
	return len(l.backend)
}

func (l *BasicLog) Entry(idx int) *Entry {
	return &l.backend[idx]
}

func (l *BasicLog) Entries(start, end int) []Entry {
	if start < 0 || end > l.Count() {
		return make([]Entry, 0)
	}

	return l.backend[start:end]
}
