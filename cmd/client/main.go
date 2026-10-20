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

var first string 
var second string
var third string

func (m model) Init() tea.Cmd{ return nil} 


func (m model) Update(msg tea.Msg)(tea.Model, tea.Cmd){
	switch msg := msg.(type){
	case tea.KeyPressMsg:
		m.lastKey = msg.String()
		switch msg.String(){
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
					
		case "down", "j":
			if m.count >0{
			m.count--}

		}}
return m, nil
	}

	func (m model) View() tea.View{
		return tea.NewView(fmt.Sprintf("%s[%s]Bread\n%s[%s]Baguette\n%s[%s]Milk",first, second,third))}

		func main(){
		if _, err := tea.NewProgram(model{}).Run(); err != nil{
fmt.Println(err)
os.Exit(1)
		}

		}

