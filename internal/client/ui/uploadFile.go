package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type uploadFile struct {
	textInput     textinput.Model
	err           error
	client        KeeperClient
	successUpload bool
	ctx           context.Context
	back          tea.Model
}

func NewUploadFileModel(ctx context.Context, client KeeperClient) uploadFile {
	ti := textinput.New()
	ti.Placeholder = "Введи путь файла"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 156

	return uploadFile{
		textInput:     ti,
		err:           nil,
		client:        client,
		successUpload: false,
		ctx:           ctx,
		back:          NewDataMenu(ctx, client),
	}
}

func (m uploadFile) Init() tea.Cmd {
	return textinput.Blink
}

func (m uploadFile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			err := m.client.UploadBinaryFile(m.ctx, m.textInput.Value(), "описание")
			if err != nil {
				m.successUpload = false
				m.err = err
				m.textInput.Reset()
				return m, cmd
			}
			m.err = nil
			m.successUpload = true
			m.textInput.Reset()
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m uploadFile) View() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"Введите путь до файла\n\n%s",
		m.textInput.View(),
	))
	sb.WriteString("\n\n (ctrl+c для выхода ctrl+b вернуться назад) \n")
	if m.successUpload {
		sb.WriteString("Файл успешно загружен \n")
	}
	if m.err != nil {
		sb.WriteString("Ошибка: ")
		sb.WriteString(errorStyle.Render(m.err.Error()))
	}
	sb.WriteString("\n")
	return sb.String()
}
