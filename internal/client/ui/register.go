package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	registerFocusedButton = focusedStyle.Render("[ Зарегистрироваться ]")
	registerBlurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Зарегистрироваться"))
)

type registerModel struct {
	focusIndex int
	inputs     []textinput.Model
	login      string
	password   string
	client     KeeperClient
	err        *error
	success    bool
}

func InitialRegisterModel(client KeeperClient) registerModel {
	m := registerModel{
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

func (m registerModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m registerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlR:
			return InitialLoginModel(m.client), nil

		case tea.KeyEnter:
			if m.focusIndex == 2 {
				err := m.client.Register(context.TODO(), m.login, m.password)
				if err != nil {
					m.err = &err
				}
				m.success = true
				return m, nil
			}
			currentInput := m.inputs[m.focusIndex]
			if loginInput == m.focusIndex {
				login := currentInput.Value()
				m.login = login
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

func (m *registerModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m registerModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &registerBlurredButton
	if m.focusIndex == len(m.inputs) {
		button = &registerFocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString("\n")
	if m.success {
		b.WriteString(helpStyle.Render("Регистрация прошла успешно"))
	}
	b.WriteString(helpStyle.Render(" (ctrl+r для входа)"))
	b.WriteString(fmt.Sprintf("Логин %s \n", m.login))
	b.WriteString(fmt.Sprintf("Пароль %s \n", m.password))
	if m.err != nil {
		b.WriteString("Ошибка регистрации \n")
		b.WriteString((*m.err).Error())
	}

	return b.String()
}
