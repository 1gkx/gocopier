package wizard

// const (
// 	suggestEscape = "(esc to quit)"
// )

// type Question struct {
// 	question     string
// 	defaultValue string
// 	choices      []string
// 	answer       string
// }

// func NewQuestion(
// 	question string,
// 	defvalue string,
// 	choices []string,
// ) Question {
// 	return Question{
// 		question:     question,
// 		defaultValue: defvalue,
// 		choices:      choices,
// 	}
// }

// type model struct {
// 	index       int
// 	questions   []Question
// 	answerField textinput.Model
// }

// func initModel(q []Question) model {
// 	ti := textinput.New()
// 	ti.Placeholder = "placeholder"
// 	ti.Focus()

// 	return model{
// 		answerField: ti,
// 		questions:   q,
// 	}
// }

// func (m model) Init() tea.Cmd {
// 	return textinput.Blink
// }

// func (m model) View() string {
// 	if m.questions[m.index].defaultValue != "" {
// 		m.answerField.Placeholder = m.questions[m.index].defaultValue
// 	}
// 	return lipgloss.JoinVertical(lipgloss.Left,
// 		m.questions[m.index].question,
// 		m.answerField.View(),
// 		suggestEscape,
// 	)
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmd tea.Cmd

// 	current := &m.questions[m.index]

// 	switch msg := msg.(type) {
// 	// case tea.WindowSizeMsg:
// 	case tea.KeyMsg:
// 		switch msg.Type {
// 		case tea.KeyCtrlC, tea.KeyEsc:
// 			return m, tea.Quit
// 		case tea.KeyEnter:
// 			current.answer = m.answerField.Value()
// 			m.answerField.SetValue("")
// 			if m.index < len(m.questions)-1 {
// 				m.index++
// 			} else {
// 				return m, tea.Quit
// 			}
// 			return m, nil
// 		}
// 	}

// 	m.answerField, cmd = m.answerField.Update(msg)
// 	return m, cmd
// }
