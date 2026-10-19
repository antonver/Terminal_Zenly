package main

import(
	"fmt"
	"os"
	tea "charm.land/bubbletea/v2"
)


type model struct{
	count int
}


func (m model) Init() tea.Cmd{ return nil} 


func (m model) Update(msg tea.Msg)(tea.Model, tea.Cmd){
	switch msg := msg.(type){
	case tea.KeyPressMsg:
		switch msg.String(){
		case "q", "cntrl + c":
			return m, tea.Quit
		case "up":
			m.count++
		case "down":
			m.count--

		}}
return m, nil
	}

	func (m model) View() tea.View{
		return tea.NewView(fmt.Sprintf("There is %d == count and to change it guess", m.count))}

		func main(){
		if _, err := tea.NewProgram(model{}).Run(); err != nil{
fmt.Println(err)
os.Exit(1)
		}

		}

