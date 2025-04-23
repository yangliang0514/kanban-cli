package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	todoList = iota
	inProgressList
	doneList
)

type Model struct {
	list     list.Model
	focusCol int
}

func New() *Model {
	return &Model{}
}

func (m *Model) initList(width, height int) {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = "To Do"
	m.list.SetItems(
		[]list.Item{
			Task{status: todo, title: "Write proposal", description: "Draft the pitch deck"},
			Task{status: todo, title: "Buy cat food", description: "Try the new salmon flavor"},
			Task{status: todo, title: "Plan weekend trip", description: "Book train tickets"},
		})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
