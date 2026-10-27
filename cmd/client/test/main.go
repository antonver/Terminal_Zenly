package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	 "charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
)

type model struct {
	spinner spinner.Model
	input textinput.Model 
	greeting string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.input.Focus, m.) 
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	
	var cmd tea.Cmd
	var cmd2 tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.input.Focused(){
		if msg.String() == "esc"{
		m.input.Blur()
		}else if msg.String() == "enter"{
		m.greeting = "Hi, " + m.input.Value()
			m.input.SetValue("")
			m.input.Blur()
		} 
		}else if msg.String() == "i"{
		return m, m.input.Focus()
		} else if msg.String() == "q"{
		return m, tea.Quit
		}
	switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
				}
	
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg) 
	
	}
	m.input, cmd2 = m.input.Update(msg)
	
	return m, tea.Batch(cmd, cmd2) 
}

func (m model) View() tea.View {
	s := "Шаг 5.x\n\n"
	s += m.spinner.View() + "\n"
	s += m.input.View() + "\n"
	s += "Name:" + m.greeting
	s += "\nctrl+c: выход\n"
	return tea.NewView(s)
}

func main() {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	text := textinput.New()
	text.Placeholder = "What is your name?"
	text.Focus()
	m := model{
		spinner: spin,
		input: text,
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
