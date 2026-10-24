package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"os"
	"strings"
	"time"
)

type model struct {
	isError        bool
	list           []string
	chosen         map[int]struct{}
	cursor         int
	lastKey        string
	windowWidth    int
	windowHeight   int
	loading        bool
	loadingCounter int
	t              int
}
type loadingMsg struct{}

type tickMsg struct{}

func tick(t time.Duration) tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}
func loadingF(t time.Duration) tea.Cmd {

	return func() tea.Msg {

		time.Sleep(t * time.Second)
		return loadingMsg{}
	}
}
func (m model) Init() tea.Cmd {
	return tea.Batch(loadingF(2), tick(time.Second))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		m.lastKey = msg.String()
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "j":
			if m.cursor < len(m.list)-1 {
				m.cursor++
			}
		case "down", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter", "tab":
			if _, ok := m.chosen[m.cursor]; ok {
				delete(m.chosen, m.cursor)
			} else {
				m.chosen[m.cursor] = struct{}{}
			}

		case "r":
			m.loadingCounter++
			if m.loadingCounter%2 == 0 {
				m.isError = true
			} else {
				m.isError = false
				m.loading = true
				return m, loadingF(2)
			}
		}
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.windowWidth = msg.Width
	case loadingMsg:
		m.loading = false

	case tickMsg:
		m.t++
		return m, tick(time.Second)
	}
	return m, nil
}
func (m model) View() tea.View {
	var b strings.Builder
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	header := style.Render("List")
	b.WriteString(header)
	b.WriteString("\n")
	if m.isError {
		b.WriteString("Error, please retry")
	} else if !m.loading {
		for number, item := range m.list {
			cursor := ""
			chosen := ""
			if m.cursor == number {
				cursor = ">"
			}
			if _, ok := m.chosen[number]; ok {
				chosen = "X"
			}
			fmt.Fprintf(&b, "%s[%s]%s\n", cursor, chosen, item)
		}
		fmt.Fprintf(&b, "move k,j to up and down, enter to chose the item\n")
	} else if m.loading {
		b.WriteString("Loading...")
	}

	fmt.Fprintf(&b, "sizes: %d x %d\n", m.windowHeight, m.windowWidth)
	b.WriteString(fmt.Sprint("Application time run in seconds: ", m.t))
	style2 := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	v := tea.NewView(style2.Render(b.String()))
	v.AltScreen = true
	return v
}
func main() {
	if _, err := tea.NewProgram(model{
		list:           []string{"bread", "milk", "jus", "crunchy sous"},
		chosen:         map[int]struct{}{},
		loading:        true,
		loadingCounter: 1,
	}).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
