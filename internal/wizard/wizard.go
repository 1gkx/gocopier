package wizard

import (
	"github.com/1gkx/gocopier/internal/configurator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Question struct {
	variableName string
	question     string
	answer       string
	input        Input
}

type Main struct {
	width     int
	height    int
	index     int
	questions []Question
	done      bool
}

func NewQ(questions []Question) *Main {
	return &Main{questions: questions}
}

func (m *Main) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

func (m Main) Init() tea.Cmd {
	return m.questions[m.index].input.Blink
}

func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := &m.questions[m.index]
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.index == len(m.questions)-1 {
				current.answer = current.input.Value()
				m.done = true
				return m, tea.Quit
			}
			current.answer = current.input.Value()
			m.Next()
			return m, current.input.Blur
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m Main) View() string {
	current := m.questions[m.index]
	if m.done {
		return ""
	}
	if m.width == 0 {
		return "loading..."
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		current.question,
		current.input.View(),
	)
}

func NewQuestion(q string) Question {
	return Question{question: q}
}

func NewInputText(vn string, q string) Question {
	qt := NewQuestion(q)
	md := NewInputField()
	qt.input = md
	qt.variableName = vn
	return qt
}

func NewChoiseText(vn string, q string, choises []string) Question {
	qt := NewQuestion(q)
	md := NewChoiceField(choises)
	qt.input = md
	qt.variableName = vn
	return qt
}

func NewQst(vn string, q configurator.Question) Question {
	if len(q.Choices) > 0 {
		return NewChoiseText(vn, q.Title, q.Choices)
	}
	return NewInputText(vn, q.Title)
}

func New(q map[string]configurator.Question) (map[string]any, error) {
	questions := make([]Question, 0, len(q))
	for k, v := range q {
		questions = append(questions, NewQst(k, v))
	}

	qs := NewQ(questions)

	p := tea.NewProgram(qs)

	_, err := p.Run()
	if err != nil {
		return nil, err
	}

	resp := make(map[string]any, len(qs.questions))
	for _, q := range qs.questions {
		resp[q.variableName] = q.answer
	}
	return resp, nil
}
