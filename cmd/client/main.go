package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"os"
	"strings"
	"time"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
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
	spin spinner.Model
	input textinput.Model
}
type loadingMsg []string 
type errMsg struct{}
type tickMsg struct{}

func tick(t time.Duration) tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}
func loadingF(count int) tea.Cmd {
	if count % 2 == 0 && count != 0 {
	return func() tea.Msg{
	time.Sleep(2 * time.Second)
	return errMsg{}
	}
	} 
return func() tea.Msg{
	time.Sleep(2 * time.Second)
	return loadingMsg{"bread", "milk", "jus", "crunchy sous"}
}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(loadingF(m.loadingCounter), tick(time.Second), m.spin.Tick, m.input.Focus())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmd2 tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.input.Focused(){
			m.input, cmd2 = m.input.Update(msg)
		if msg.String() == "enter"{
			if m.input.Value() == "" || m.input.Value() == " "{
			
			m.input.Blur()
			}else{
			m.list = append(m.list, m.input.Value())	
			
			m.input.SetValue("")
			m.input.Blur()
		}}
			
			return m, cmd2
		}else{
		
		
		m.lastKey = msg.String()
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "down", "j":
			if m.cursor < len(m.list)-1 {
				m.cursor++
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "tab":
			if _, ok := m.chosen[m.cursor]; ok {
				delete(m.chosen, m.cursor)
			} else {
				m.chosen[m.cursor] = struct{}{}
			}
		case "a":
			return m, m.input.Focus() 
		case "r":
			m.loadingCounter++
			m.isError = false
				m.loading = true
				return m, tea.Batch(loadingF(m.loadingCounter), m.spin.Tick)	
			}
		}
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.windowWidth = msg.Width
	case loadingMsg:
		m.loading = false
		m.list = msg
	case tickMsg:
		m.t++
		return m, tick(time.Second)
	
	case errMsg:
		m.loading = false
		m.isError = true
	case spinner.TickMsg:
		m.spin, cmd = m.spin.Update(msg)
		return m, cmd
	}
	m.input, cmd2 = m.input.Update(msg)
	return m, cmd2 
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
		b.WriteString(m.input.View())
		fmt.Fprintf(&b, "move k,j to up and down, enter to chose the item\n")
	} else if m.loading {
		b.WriteString(m.spin.View())
	}

	fmt.Fprintf(&b, "sizes: %d x %d\n", m.windowHeight, m.windowWidth)
	b.WriteString(fmt.Sprint("Application time run in seconds: ", m.t))
	style2 := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	v := tea.NewView(style2.Render(b.String()))
	v.AltScreen = true
	return v
}
func main() {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	text := textinput.New()
	if _, err := tea.NewProgram(model{
		chosen:         map[int]struct{}{},
		loading:        true,
		loadingCounter: 1,
		spin: spin,
		input: text,
	}).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
