package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	isError bool
	list    []string
	chosen  map[int]struct{}
	cursor  int
	// УДАЛЕНО: lastKey string — значение записывалось, но нигде не использовалось
	windowWidth    int
	windowHeight   int
	loading        bool
	loadingCounter int
	t              int
	spin           spinner.Model
	input          textinput.Model
}

type loadingMsg []string
type errMsg struct{}
type tickMsg struct{}

// ИЗМЕНЕНО: параметр t теперь используется (раньше внутри стояло time.Second).
// Параметр внутренней функции без имени, чтобы он не перекрывал внешний t.
func tick(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func loadingF(count int) tea.Cmd {
	if count%2 == 0 && count != 0 {
		return func() tea.Msg {
			time.Sleep(2 * time.Second)
			return errMsg{}
		}
	}
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		return loadingMsg{"bread", "milk", "jus", "crunchy sous"}
	}
}

func (m model) Init() tea.Cmd {
	// ИЗМЕНЕНО: убран m.input.Focus().
	// Init не может менять модель: Focus() менял копию, и она выбрасывалась.
	// Плюс по заданию программа стартует в режиме навигации, фокус даётся по "a".
	return tea.Batch(loadingF(m.loadingCounter), tick(time.Second), m.spin.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// УДАЛЕНО: var cmd2 tea.Cmd — теперь в каждой ветке одна команда, хватает одной переменной

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// ИЗМЕНЕНО: ctrl+c вынесен наверх, работает в обоих режимах
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// ===== РЕЖИМ ВВОДА =====
		if m.input.Focused() {
			// ИЗМЕНЕНО: сначала проверяем служебные клавиши,
			// и только потом отдаём остальное полю.
			// Раньше клавиша сначала уходила в поле, а потом проверялся enter.
			switch msg.String() {
			case "enter":
				// ИЗМЕНЕНО: TrimSpace вместо сравнения с "" и " ".
				// Теперь "   " (несколько пробелов) тоже считается пустой строкой.
				text := strings.TrimSpace(m.input.Value())
				if text != "" {
					m.list = append(m.list, text)
				}
				// ИЗМЕНЕНО: поле очищается всегда, а не только при пустой строке.
				// Раньше после добавления "eggs" в поле так и оставалось "eggs".
				m.input.SetValue("")
				// ИЗМЕНЕНО: Blur() только здесь и в esc.
				// ГЛАВНЫЙ БАГ: раньше Blur() вызывался после КАЖДОЙ клавиши,
				// поэтому можно было напечатать только одну букву.
				m.input.Blur()
				return m, nil

			case "esc":
				// ДОБАВЛЕНО: отмена ввода без добавления
				m.input.SetValue("")
				m.input.Blur()
				return m, nil
			}

			// Все остальные клавиши, включая q, j, k, r, a, это текст для поля.
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		// ===== РЕЖИМ НАВИГАЦИИ =====
		// УДАЛЕНО: m.lastKey = msg.String()
		// ИЗМЕНЕНО: убран лишний else, вместо него ранний return в ветке выше
		switch msg.String() {
		case "q":
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
			// ДОБАВЛЕНО: ввод только когда список показан.
			// Раньше можно было печатать во время загрузки вслепую (поле не рисовалось).
			if m.loading || m.isError {
				return m, nil
			}
			return m, m.input.Focus()
		case "r":
			// ДОБАВЛЕНО: защита от второй загрузки поверх первой
			if m.loading {
				return m, nil
			}
			m.loadingCounter++
			m.isError = false
			m.loading = true
			// ИЗМЕНЕНО: убран m.spin.Tick.
			// Спиннер крутится всё время, потому что TickMsg всегда пробрасывается ему.
			return m, loadingF(m.loadingCounter)
		}
		// ДОБАВЛЕНО: клавиши навигации не отдаём в поле, оно сейчас не в фокусе
		return m, nil

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

	// Всё остальное (например мигание курсора) отдаём полю
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	var b strings.Builder
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	header := style.Render("List")
	b.WriteString(header)
	b.WriteString("\n")

	// ИЗМЕНЕНО: порядок веток стал error → loading → список,
	// и последняя ветка теперь просто else, без повторной проверки m.loading
	if m.isError {
		// ИЗМЕНЕНО: добавлен \n и подсказка
		b.WriteString("Error, please retry (r)\n")
	} else if m.loading {
		// ИЗМЕНЕНО: добавлен текст и \n
		b.WriteString(m.spin.View() + " Loading...\n")
	} else {
		for number, item := range m.list {
			// ИЗМЕНЕНО: пробел вместо "", чтобы все строки были одной ширины
			cursor := " "
			chosen := " "
			if m.cursor == number {
				cursor = ">"
			}
			if _, ok := m.chosen[number]; ok {
				chosen = "X"
			}
			fmt.Fprintf(&b, "%s[%s]%s\n", cursor, chosen, item)
		}

		// ИЗМЕНЕНО: поле рисуется только в режиме ввода,
		// и подсказка по клавишам своя для каждого режима
		if m.input.Focused() {
			b.WriteString(m.input.View() + "\n")
			b.WriteString("enter: add | esc: cancel\n")
		} else {
			b.WriteString("k/j: move | tab: choose | a: add | r: reload | q: quit\n")
		}
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
	// ДОБАВЛЕНО: placeholder по заданию
	text.Placeholder = "New item"
	// ДОБАВЛЕНО: ширина, иначе placeholder обрезается до одной буквы.
	// Если не компилируется, проверь имя метода:
	// go doc charm.land/bubbles/v2/textinput Model
	text.SetWidth(30)

	if _, err := tea.NewProgram(model{
		chosen:         map[int]struct{}{},
		loading:        true,
		loadingCounter: 1,
		spin:           spin,
		input:          text,
	}).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
