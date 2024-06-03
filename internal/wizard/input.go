package wizard

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Blink() tea.Msg
	Blur() tea.Msg
	Focus() tea.Cmd
	SetValue(string)
	Value() string
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
}

type inputAnswerField struct {
	textinput textinput.Model
}

func NewInputField() *inputAnswerField {
	a := inputAnswerField{}

	ti := textinput.New()
	ti.Placeholder = "Your answer here"
	ti.Focus()

	a.textinput = ti
	return &a
}

func (a *inputAnswerField) Blink() tea.Msg {
	return textinput.Blink()
}

func (a *inputAnswerField) Init() tea.Cmd {
	return nil
}

func (a *inputAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	a.textinput, cmd = a.textinput.Update(msg)
	return a, cmd
}

func (a *inputAnswerField) View() string {
	return a.textinput.View()
}

func (a *inputAnswerField) Focus() tea.Cmd {
	return a.textinput.Focus()
}

func (a *inputAnswerField) SetValue(s string) {
	a.textinput.SetValue(s)
}

func (a *inputAnswerField) Blur() tea.Msg {
	return a.textinput.Blur
}

func (a *inputAnswerField) Value() string {
	return a.textinput.Value()
}

/*** choices ***/
type choiceAnswerField struct {
	cursor  int
	choice  string
	choices []string
}

func NewChoiceField(choices []string) *choiceAnswerField {
	a := choiceAnswerField{
		choices: choices,
	}
	return &a
}

func (a *choiceAnswerField) Blink() tea.Msg {
	return textinput.Blink()
}

func (a *choiceAnswerField) Init() tea.Cmd {
	return nil
}

func (a *choiceAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			a.cursor++
			if a.cursor >= len(a.choices) {
				a.cursor = 0
			}
			a.choice = a.choices[a.cursor]

		case tea.KeyUp:
			a.cursor--
			if a.cursor < 0 {
				a.cursor = len(a.choices) - 1
			}
			a.choice = a.choices[a.cursor]
		}
	}

	a.choice = a.choices[a.cursor]
	return a, nil
}

func (a *choiceAnswerField) View() string {
	s := strings.Builder{}

	for i := 0; i < len(a.choices); i++ {
		if a.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(a.choices[i])
		s.WriteString("\n")
	}

	return s.String()
}

func (a *choiceAnswerField) Focus() tea.Cmd {
	return nil
}

func (a *choiceAnswerField) SetValue(_ string) {
}

func (a *choiceAnswerField) Blur() tea.Msg {
	return nil
}

func (a *choiceAnswerField) Value() string {
	return a.choice
}
