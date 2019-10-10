package log

import "cjvirtucio87/distributed-todo-go/internal/dto"

type Log interface {
	AddEntries(entries []dto.Entry)
	Count() int
	Entry(idx int) (dto.Entry, bool)
}
