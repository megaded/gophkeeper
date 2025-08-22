package ui

import (
	"context"
	"errors"
	"fmt"
	"gophkeeper/internal/server/dto"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

const (
	ccn = iota
	exp
	cvv
)

type modeType int

const (
	new modeType = iota
	edit
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle        = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle     = lipgloss.NewStyle().Foreground(darkGray)
	saveFocusedButton = focusedStyle.Render("[ Сохранить ]")
	saveBlurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Сохранить"))
)

type creditCardModel struct {
	inputs  []textinput.Model
	focused int
	err     error
	ccn     string
	exp     string
	cvv     string
	mode    modeType
	client  KeeperClient
}

// Validator functions to ensure valid input
func ccnValidator(s string) error {
	// Credit Card Number should a string less than 20 digits
	// It should include 16 integers and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("CCN is too long")
	}

	if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("CCN is invalid")
	}

	// The last digit should be a number unless it is a multiple of 4 in which
	// case it should be a space
	if len(s)%5 == 0 && s[len(s)-1] != ' ' {
		return fmt.Errorf("CCN must separate groups with spaces")
	}

	// The remaining digits should be integers
	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func expValidator(s string) error {
	// The 3 character should be a slash (/)
	// The rest should be numbers
	e := strings.ReplaceAll(s, "/", "")
	_, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		return fmt.Errorf("EXP is invalid")
	}

	// There should be only one slash and it should be in the 2nd index (3rd character)
	if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
		return fmt.Errorf("EXP is invalid")
	}

	return nil
}

func cvvValidator(s string) error {
	_, err := strconv.ParseInt(s, 10, 64)
	if !validateCardNumber(s) {
		return errors.New("Invalid card number")
	}
	return err
}

func InitCreditCardModel(client KeeperClient) creditCardModel {
	return creditCardModel{
		inputs:  getInputs(),
		focused: 0,
		err:     nil,
		mode:    new,
		client:  client,
	}
}

func InitialCreditCardEditModel(client KeeperClient, ccn string, exp string, cvv string) creditCardModel {
	return creditCardModel{
		inputs:  getInputs(),
		focused: 0,
		err:     nil,
		mode:    edit,
		ccn:     ccn,
		exp:     exp,
		cvv:     cvv,
		client:  client,
	}
}

func getInputs() []textinput.Model {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	inputs[ccn] = textinput.New()
	inputs[ccn].Placeholder = "xxxx xxxx xxxx xxxx"
	inputs[ccn].Focus()
	inputs[ccn].CharLimit = 19
	inputs[ccn].Width = 30
	inputs[ccn].Validate = ccnValidator

	inputs[exp] = textinput.New()
	inputs[exp].Placeholder = "MM/YY "
	inputs[exp].CharLimit = 5
	inputs[exp].Width = 5
	inputs[exp].Validate = expValidator

	inputs[cvv] = textinput.New()
	inputs[cvv].Placeholder = "XXX"
	inputs[cvv].CharLimit = 3
	inputs[cvv].Width = 5
	inputs[cvv].Validate = cvvValidator
	return inputs
}

func (m creditCardModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m creditCardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs) {
				err := m.client.AddCreditCard(context.Background(), dto.Card{Number: m.inputs[ccn].Value(), Exp: m.inputs[exp].Value(), CVV: m.inputs[cvv].Value()})
				if err != nil {
					fmt.Println(err.Error())
					return m, tea.Quit
				}
				return NewDataMenu(m.client), tea.Quit
			}
			m.focused++
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		if m.focused <= len(m.inputs)-1 {
			m.inputs[m.focused].Focus()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		input := m.inputs[i]
		input, cmds[i] = m.inputs[i].Update(msg)
		if i == ccn {
			input.SetValue(formatCCNInput(input.Value()))
			input.CursorEnd()
		}
		m.inputs[i] = input
	}
	return m, tea.Batch(cmds...)
}

func formatCCNInput(raw string) string {
	rawLen := len(raw)
	if rawLen == 4 || rawLen == 9 || rawLen == 14 {
		raw = raw + " "
	}
	return raw
}

func (m creditCardModel) View() string {
	var title string
	if m.mode == new {
		title = "Новая кредитная карта"
	} else {
		title = "Редактирование кредитной карты"
	}
	button := saveBlurredButton
	if m.focused == len(m.inputs) {
		button = saveFocusedButton
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		` %s

 %s
 %s

 %s  %s
 %s  %s

 %s
`, title,
		inputStyle.Width(30).Render("Card Number"),
		m.inputs[ccn].View(),
		inputStyle.Width(6).Render("EXP"),
		inputStyle.Width(6).Render("CVV"),
		m.inputs[exp].View(),
		m.inputs[cvv].View(),
		button,
	) + "\n")
	if m.inputs != nil {
		err := m.inputs[ccn].Err
		if m.inputs[ccn].Err != nil {
			sb.WriteString(err.Error())
		}

	}
	return sb.String()
}

// prevInput focuses the previous input field
func (m *creditCardModel) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func validateCardNumber(number string) bool {
	var sum int
	double := false

	for i := len(number) - 1; i >= 0; i-- {
		r := number[i]

		if r < '0' || r > '9' {
			return false
		}

		digit := int(r - '0')

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	return sum%10 == 0
}
