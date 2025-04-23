package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	todoList = iota
	inProgressList
	doneList
)

type Model struct {
	list     []list.Model
	focusCol int
}

func New() *Model {
	return &Model{}
}

func (m *Model) initList(width, height int) {
	m.list = make([]list.Model, 3)

	m.list[todoList] = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list[todoList].Title = "To Do"
	m.list[todoList].SetItems(todoMockData())

	m.list[inProgressList] = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list[inProgressList].Title = "In Progress"
	m.list[inProgressList].SetItems(inProgressMockData())

	m.list[doneList] = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list[doneList].Title = "Done"
	m.list[doneList].SetItems(doneMockData())
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.handleKeyPress(msg.String())
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list[m.focusCol], cmd = m.list[m.focusCol].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if len(m.list) == 0 {
		return "Loading..."
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		m.list[todoList].View(),
		m.list[inProgressList].View(),
		m.list[doneList].View(),
	)
}

func (m *Model) handleKeyPress(msgStr string) {
	switch msgStr {
	case tea.KeyRight.String():
		if m.focusCol >= doneList {
			m.focusCol = todoList
		} else {
			m.focusCol++
		}
	case tea.KeyLeft.String():
		if m.focusCol <= todoList {
			m.focusCol = doneList
		} else {
			m.focusCol--
		}
	}
}

// temporary mock data
func todoMockData() []list.Item {
	return []list.Item{
		Task{status: todo, title: "Write proposal", description: "Draft the pitch deck"},
		Task{status: todo, title: "Buy cat food", description: "Try the new salmon flavor"},
		Task{status: todo, title: "Plan weekend trip", description: "Book train tickets"},
	}
}

func inProgressMockData() []list.Item {
	return []list.Item{
		&Task{status: inProgress, title: "Fix login bug", description: "Investigating session timeout issue"},
		&Task{status: inProgress, title: "Refactor profile page", description: "Split into smaller components"},
		&Task{status: inProgress, title: "Write tests", description: "Add coverage for user service"},
	}
}

func doneMockData() []list.Item {
	return []list.Item{
		&Task{status: done, title: "Set up CI", description: "Configured GitHub Actions for builds"},
		&Task{status: done, title: "Deploy v1.0", description: "Released first version to production"},
		&Task{status: done, title: "Onboard new dev", description: "Walked them through codebase"},
	}
}
