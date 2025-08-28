package ui

// A simple program demonstrating the textarea component from the Bubbles
// component library.

import (
	"context"
	"fmt"
	"gophkeeper/internal/server/dto"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type textModel struct {
	textarea textarea.Model
	err      error
	mode     modeType
	content  string
	ctx      context.Context
	client   KeeperClient
	id       uint
}

func NewTextModel(ctx context.Context, client KeeperClient) textModel {
	ti := textarea.New()
	ti.Placeholder = "Введите текст"
	ti.Focus()

	return textModel{
		textarea: ti,
		err:      nil,
		mode:     new,
		ctx:      ctx,
		client:   client,
	}
}

func NewEditTextModel(ctx context.Context, client KeeperClient, content string, id uint) textModel {
	ti := textarea.New()
	ti.SetValue(content)
	ti.Focus()

	return textModel{
		textarea: ti,
		err:      nil,
		mode:     edit,
		ctx:      ctx,
		client:   client,
		content:  content,
		id:       id,
	}
}

func (m textModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlS:
			if m.mode == new {
				err := m.client.UploadText(context.Background(), m.textarea.Value(), "")
				if err != nil {
					fmt.Println(err.Error())
					return m, tea.Quit
				}
			}
			if m.mode == edit {
				err := m.client.UpdateText(context.Background(), dto.Text{Content: m.textarea.Value(), Id: m.id})
				if err != nil {
					fmt.Println(err.Error())
					return m, tea.Quit
				}
			}

			return NewDataMenu(m.ctx, m.client), nil
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m textModel) View() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"Новые текстовые данные.\n\n\n%s",
		m.textarea.View(),
	))
	sb.WriteString("\n\n")
	sb.WriteString("\n\n (ctrl+c для выхода | ctrl+b вернуться назад | ctrl+s сохранить) \n")
	return sb.String()
}
