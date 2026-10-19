package main

import(
	"fmt"
	"os"
	tea "charm.land/bubbletea/v2"
)


type model struct{
	count int
	lastKey string
}


func (m model) Init() tea.Cmd{ return nil} 


func (m model) Update(msg tea.Msg)(tea.Model, tea.Cmd){
	switch msg := msg.(type){
	case tea.KeyPressMsg:
		m.lastKey = msg.String()
		switch msg.String(){
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
	if m.count < 10{	
			m.count++}
		case "down":
			if m.count >0{
			m.count--}

		}}
return m, nil
	}

	func (m model) View() tea.View{
		return tea.NewView(fmt.Sprintf("There is %d == count and to change it guess, the word: %s", m.count, m.lastKey))}

		func main(){
		if _, err := tea.NewProgram(model{}).Run(); err != nil{
fmt.Println(err)
os.Exit(1)
		}

		}

