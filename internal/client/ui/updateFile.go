package ui

import (
	"context"
	"fmt"
	fileutil "gophkeeper/internal/client/proto/fileUtil"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type updateFile struct {
	textInput       textinput.Model
	err             error
	client          KeeperClient
	successUpload   bool
	ctx             context.Context
	back            tea.Model
	id              uint
	fileName        string
	successDownload bool
}

func NewUpdateFileModel(ctx context.Context, client KeeperClient, id uint, fileName string) updateFile {
	ti := textinput.New()
	ti.Placeholder = "Введи путь файла"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 156

	return updateFile{
		textInput:     ti,
		err:           nil,
		client:        client,
		successUpload: false,
		ctx:           ctx,
		back:          NewDataMenu(ctx, client),
		id:            id,
		fileName:      fileName,
	}
}

func (m updateFile) Init() tea.Cmd {
	return textinput.Blink
}

func (m updateFile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			err := m.client.UpdateBinaryFile(m.ctx, m.textInput.Value(), "описание", m.id)
			if err != nil {
				m.successUpload = false
				m.err = err
				m.textInput.Reset()
				return m, cmd
			}
			m.err = nil
			m.successUpload = true
			m.successDownload = false
			m.textInput.Reset()
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlT:
			err := m.client.DownloadBinaryFile(m.ctx, m.id)
			if err != nil {
				m.successDownload = false
				m.err = err
				m.textInput.Reset()
				return m, cmd
			}
			m.err = nil
			m.successUpload = false
			m.successDownload = true
			m.textInput.Reset()
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m updateFile) View() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"Файл %s\n\n%s",
		m.fileName,
		m.textInput.View(),
	))
	sb.WriteString("\n\n (ctrl+c для выхода ctrl+b вернуться назад ctrl + t загрузить файл) \n")
	if m.successUpload {
		sb.WriteString("Файл успешно загружен \n")
	}
	if m.successDownload {
		sb.WriteString(fmt.Sprintf("Файл успешно скачен путь %s\n", fileutil.GetDownloadDir(m.fileName)))
	}
	if m.err != nil {
		sb.WriteString("Ошибка: ")
		sb.WriteString(errorStyle.Render(m.err.Error()))
	}
	sb.WriteString("\n")
	return sb.String()
}
