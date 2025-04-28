package main

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	listStyle  = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("0")).AlignVertical(lipgloss.Center)
	focusStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")).AlignVertical(lipgloss.Center)
	modalStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("10")).Align(lipgloss.Center, lipgloss.Center)
)

const (
	listsView = iota
	moveToView
	editView
)

type Model struct {
	currentView int
	lists       []list.Model
	focusList   Status
	width       int
	height      int
	quitting    bool
	moveToModel MoveToModel
	undoStack   []UndoItem
}

type MoveToModel struct {
	optionList  list.Model
	toList      Status
	initialized bool
}

type UndoItem struct {
	list  Status
	item  list.Item
	index int
}

func New() *Model {
	return &Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msgStr := msg.String(); msgStr == tea.KeyCtrlC.String() || msgStr == "q" {
			m.quitting = true
			return m, tea.Quit
		}

		return m.handleKeyStroke(msg)
	case tea.WindowSizeMsg:
		if len(m.lists) == 0 {
			m.initList(msg.Width, msg.Height)
		}
	}

	return m, nil
}

func (m Model) View() string {
	if len(m.lists) == 0 || m.quitting {
		return ""
	}

	switch m.currentView {
	case listsView:
		return m.renderListsView()
	case moveToView:
		return m.renderMoveToView()
	}

	return ""
}

func (m *Model) initList(width, height int) {
	m.lists = make([]list.Model, 3)
	m.width = width
	m.height = height

	m.lists[todo] = list.New(todoMockData(), list.NewDefaultDelegate(), m.width/4, m.height-5)
	m.lists[todo].Title = "To Do"

	m.lists[inProgress] = list.New(inProgressMockData(), list.NewDefaultDelegate(), m.width/4, m.height-5)
	m.lists[inProgress].Title = "In Progress"

	m.lists[done] = list.New(doneMockData(), list.NewDefaultDelegate(), m.width/4, m.height-5)
	m.lists[done].Title = "Done"

	m.lists[todo].SetShowHelp(false)
	m.lists[inProgress].SetShowHelp(false)
	m.lists[done].SetShowHelp(false)
}

func (m *Model) initMoveToModel() {
	options := []list.Item{
		ListOption{list: todo, title: "Todo"},
		ListOption{list: inProgress, title: "In Progress"},
		ListOption{list: done, title: "Done"},
	}

	list := list.New(options, list.NewDefaultDelegate(), m.width/2, m.height/2)
	list.SetShowHelp(false)

	m.moveToModel = MoveToModel{optionList: list, initialized: true}
}

func (m *Model) handleKeyStroke(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.currentView {
	case listsView:
		m.handleKeyListView(msg.String())

		var cmd tea.Cmd
		m.lists[m.focusList], cmd = m.lists[m.focusList].Update(msg)
		return m, cmd
	case moveToView:
		cmd1 := m.handleKeyMoveToView(msg.String())

		var cmd2 tea.Cmd
		m.moveToModel.optionList, cmd2 = m.moveToModel.optionList.Update(msg)
		return m, tea.Batch(cmd1, cmd2)
	}

	return m, nil
}

func (m *Model) handleKeyListView(msgStr string) {
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
	case "d":
		currentList := &m.lists[m.focusList]
		currentItem := currentList.SelectedItem()
		currentIndex := currentList.Index()

		if currentItem == nil {
			return
		}

		items := currentList.Items()
		items = slices.Delete(items, currentIndex, currentIndex+1)
		currentList.SetItems(items)
		m.pushToUndoStack(m.focusList, currentItem, currentIndex)
	case "u":
		targetList, item, targetIndex := m.popFromUndoStack()

		if item == nil {
			return
		}

		undoItemList := &m.lists[targetList]
		items := undoItemList.Items()
		items = slices.Insert(items, targetIndex, item)
		undoItemList.SetItems(items)
		undoItemList.Select(targetIndex)
		m.focusList = targetList
	case tea.KeyEnter.String():
		m.currentView = moveToView

		if !m.moveToModel.initialized {
			m.initMoveToModel()
		}

	}
}

func (m *Model) handleKeyMoveToView(msgStr string) tea.Cmd {
	switch msgStr {
	case "b":
		m.currentView = listsView
	case tea.KeyEnter.String():
		selectedOption := m.moveToModel.optionList.SelectedItem().(ListOption)

		m.currentView = listsView

		return m.moveTask(int(selectedOption.list))
	}

	return nil
}

func (m *Model) moveTask(toList int) tea.Cmd {
	if int(m.focusList) == toList {
		return nil
	}

	currentList := m.lists[m.focusList]
	newTargetListItems := append(m.lists[toList].Items(), currentList.SelectedItem())

	m.lists[m.focusList].RemoveItem(currentList.Index())
	cmd := m.lists[toList].SetItems(newTargetListItems)
	m.focusList = Status(toList)
	m.lists[m.focusList].Select(len(newTargetListItems) - 1)

	return cmd
}

func (m Model) renderListWithStyle(list Status) string {
	if list == m.focusList {
		return focusStyle.Render(m.lists[list].View())
	}
	return listStyle.Render(m.lists[list].View())
}

func (m Model) renderListsView() string {
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			m.renderListWithStyle(todo),
			m.renderListWithStyle(inProgress),
			m.renderListWithStyle(done),
		),
	)
}

func (m *Model) renderMoveToView() string {
	selectedItem := m.lists[m.focusList].SelectedItem().(Task)

	m.moveToModel.optionList.Title = fmt.Sprintf("Move task [ %s ] to...", selectedItem.Title())

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modalStyle.Render(m.moveToModel.optionList.View()))
}

func (m *Model) pushToUndoStack(list Status, item list.Item, index int) {
	m.undoStack = append(m.undoStack, UndoItem{list: list, item: item, index: index})
}

func (m *Model) popFromUndoStack() (Status, list.Item, int) {
	if len(m.undoStack) == 0 {
		return 0, nil, 0
	}

	lastItem := m.undoStack[len(m.undoStack)-1]
	m.undoStack = m.undoStack[:len(m.undoStack)-1]

	return lastItem.list, lastItem.item, lastItem.index
}

// temporary mock data
func todoMockData() []list.Item {
	return []list.Item{
		Task{title: "Write proposal", description: "Draft the pitch deck"},
		Task{title: "Buy cat food", description: "Try the new salmon flavor"},
		Task{title: "Plan weekend trip", description: "Book train tickets"},
	}
}

func inProgressMockData() []list.Item {
	return []list.Item{
		Task{title: "Fix login bug", description: "Investigating session timeout issue"},
		Task{title: "Refactor profile page", description: "Split into smaller components"},
		Task{title: "Write tests", description: "Add coverage for user service"},
		Task{title: "Deploy v1.0", description: "Released first version to production"},
	}
}

func doneMockData() []list.Item {
	return []list.Item{
		Task{title: "Set up CI", description: "Configured GitHub Actions for builds"},
		Task{title: "Onboard new dev", description: "Walked them through codebase"},
	}
}
