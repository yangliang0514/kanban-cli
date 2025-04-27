package main

type Status int

const (
	todo Status = iota
	inProgress
	done
)

// this struct implements list.Item interface
type Task struct {
	title       string
	description string
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

func (t Task) FilterValue() string {
	return t.title
}

type ListOption struct {
	list  Status
	title string
}

func (l ListOption) Title() string {
	return l.title
}

func (l ListOption) Description() string {
	return ""
}

func (l ListOption) FilterValue() string {
	return l.title
}
