package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	listStyle  = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("0")).AlignVertical(lipgloss.Center)
	focusStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")).AlignVertical(lipgloss.Center)
)

type Model struct {
	lists     []list.Model
	focusList Status
	width     int
	height    int
}

func New() *Model {
	return &Model{}
}

func (m *Model) initList(width, height int) {
	m.lists = make([]list.Model, 3)
	m.width = width
	m.height = height

	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), m.width/4, m.height-5)

	m.lists[todo] = defaultList
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems(todoMockData())

	m.lists[inProgress] = defaultList
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems(inProgressMockData())

	m.lists[done] = defaultList
	m.lists[done].Title = "Done"
	m.lists[done].SetItems(doneMockData())

	m.lists[todo].SetShowHelp(false)
	m.lists[inProgress].SetShowHelp(false)
	m.lists[done].SetShowHelp(false)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.handleKeyPress(msg.String())
	case tea.WindowSizeMsg:
		if len(m.lists) == 0 {
			m.initList(msg.Width, msg.Height)
		}
	}

	var cmd tea.Cmd
	m.lists[m.focusList], cmd = m.lists[m.focusList].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if len(m.lists) == 0 {
		return "Loading..."
	}

	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			m.renderListWithStyle(todo),
			m.renderListWithStyle(inProgress),
			m.renderListWithStyle(done),
		),
	)

}

func (m *Model) handleKeyPress(msgStr string) {
	switch msgStr {
	case tea.KeyRight.String(), "l":
		if m.focusList >= done {
			m.focusList = todo
		} else {
			m.focusList++
		}
	case tea.KeyLeft.String(), "h":
		if m.focusList <= todo {
			m.focusList = done
		} else {
			m.focusList--
		}
	}
}

func (m Model) renderListWithStyle(list Status) string {
	if list == m.focusList {
		return focusStyle.Render(m.lists[list].View())
	}
	return listStyle.Render(m.lists[list].View())
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
		&Task{status: inProgress, title: "Deploy v1.0", description: "Released first version to production"},
	}
}

func doneMockData() []list.Item {
	return []list.Item{
		&Task{status: done, title: "Set up CI", description: "Configured GitHub Actions for builds"},
		&Task{status: done, title: "Onboard new dev", description: "Walked them through codebase"},
	}
}
