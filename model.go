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
	list []list.Model
	err  error
}

func initialModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg.String())
	}
	return m, nil
}

func (m Model) View() string {
	return "Hello there"
}

func (m Model) handleKeyPress(msgStr string) (tea.Model, tea.Cmd) {
	switch msgStr {
	case "ctrl+c", "q":
		return m, tea.Quit
	}

	return m, nil
}
