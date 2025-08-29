package ui

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"context"
	"fmt"
	"gophkeeper/internal/dto"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type creditialModel struct {
	focusIndex  int
	inputs      []textinput.Model
	login       string
	password    string
	description string
	client      KeeperClient
	err         error
	ctx         context.Context
	cred        credential
	mode        modeType
}

func NewCreditialModel(ctx context.Context, client KeeperClient) creditialModel {
	m := creditialModel{
		ctx:    ctx,
		client: client,
		inputs: make([]textinput.Model, 2),
		mode:   new,
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

func InitialEditCreditialModel(ctx context.Context, client KeeperClient, cred credential) creditialModel {
	m := creditialModel{
		ctx:    ctx,
		client: client,
		inputs: make([]textinput.Model, 2),
		mode:   edit,
		cred:   cred,
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
			t.SetValue(cred.login)
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Prompt = "Password :"
			t.SetValue(cred.password)
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m creditialModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m creditialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEnter:
			if m.focusIndex == len(m.inputs) {
				if m.mode == edit {
					err := m.client.UpdateCreditial(context.Background(), dto.Credentials{Id: m.cred.id})
					if err != nil {
						fmt.Println(err.Error())
						return m, tea.Quit
					}
				}
				if m.mode == new {
					err := m.client.AddCredentials(context.Background(), dto.Credentials{Login: m.inputs[0].Value(), Password: m.inputs[1].Value()})
					if err != nil {
						fmt.Println(err.Error())
						return m, tea.Quit
					}
				}
				return NewDataMenu(m.ctx, m.client), nil
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

					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}

				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)

		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *creditialModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m creditialModel) View() string {
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
	return b.String()
}
