package ui

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	errorStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Background(lipgloss.Color("235"))
	loginFocusedButton = focusedStyle.Render("[ Войти ]")
	loginBlurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Войти"))
	loginInput         = 0
	passwordInput      = 1
)

type loginModel struct {
	focusIndex int
	inputs     []textinput.Model
	login      string
	password   string
	token      string
	client     KeeperClient
	err        error
}

func InitialLoginModel(client KeeperClient) loginModel {
	m := loginModel{
		client: client,
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		t.Cursor.SetMode(cursor.CursorBlink)

		switch i {
		case 0:
			t.Prompt = "Login: "
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Prompt = "Password :"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlR:

			return InitialRegisterModel(m.client), nil

		case tea.KeyEnter:
			s := msg.String()
			if s == "enter" && m.focusIndex == len(m.inputs) {
				r, err := m.client.Login(context.Background(), m.login, m.password)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.token = r
				return m, nil
			}
			currentInput := m.inputs[m.focusIndex]
			if loginInput == m.focusIndex {
				m.login = currentInput.Value()
			}

			if passwordInput == m.focusIndex {
				m.password = currentInput.Value()
			}
			m.focusIndex++

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)

		}
	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *loginModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m loginModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &loginBlurredButton
	if m.focusIndex == len(m.inputs) {
		button = &loginFocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render(" (ctrl+r для регистрации)"))
	b.WriteString("\n")
	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Ошибка %s", m.err.Error())))
	}
	if m.token != "" {
		b.WriteString("\n")
		b.WriteString(m.token)
	}
	return b.String()
}
