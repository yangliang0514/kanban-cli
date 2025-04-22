package main

type Status int

const (
	todo Status = iota
	inprogress
	done
)

type Task struct {
	status      Status
	title       string
	description string
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) Description() string {
	return t.description
}

func (t *Task) filterValue() string {
	return t.title
}
