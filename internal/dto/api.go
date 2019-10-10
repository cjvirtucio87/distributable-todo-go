package dto

type EntryInfo struct {
	Entries   []Entry
	NextIndex int
}

type Entry struct {
	Id      int
	Command string
}

type Message struct {
	Entries []Entry
}
