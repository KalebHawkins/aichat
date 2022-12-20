package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/KalebHawkins/ggpt3"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error
type respMsg string

type keyMap struct {
	Submit      key.Binding
	ClearPrompt key.Binding
	Quit        key.Binding
}

var keys = keyMap{
	Submit: key.NewBinding(
		key.WithKeys(tea.KeyCtrlS.String()),
		key.WithHelp(tea.KeyCtrlS.String(), "submit query"),
	),
	ClearPrompt: key.NewBinding(
		key.WithKeys(tea.KeyCtrlC.String()),
		key.WithHelp(tea.KeyCtrlC.String(), "clear prompt"),
	),
	Quit: key.NewBinding(
		key.WithKeys(tea.KeyEsc.String()),
		key.WithHelp(tea.KeyEsc.String(), "quit"),
	),
}

const (
	divisor = 4
)

var (
	titleStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Align(lipgloss.Center)
	textAreaStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	helpStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	versionInfoStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Height(3)
)

// Used for versioning information.
var (
	Version string
	Commit  string
)

type Model struct {
	title    string
	client   ggpt3.Client
	textArea textarea.Model
	keys     keyMap
	help     []string
	quitting bool
}

func newModel() Model {
	ta := textarea.New()
	ta.Placeholder = "You can ask me whatever you want."
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.Focus()

	help := []string{
		"Submit Query - Ctrl+S",
		"Clear Text   - Ctrl+C",
		"Quit         - Esc",
	}

	return Model{
		title:    "---Welcome to AIChat---\nPowered by OpenAI",
		client:   *ggpt3.NewClient(os.Getenv("AI_CHAT_KEY")),
		textArea: ta,
		keys:     keys,
		help:     help,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w := msg.Width - msg.Width/divisor
		h := msg.Height - msg.Height/divisor

		titleStyle.Width(w - titleStyle.GetHorizontalFrameSize())
		textAreaStyle.Width(w - textAreaStyle.GetHorizontalFrameSize())
		textAreaStyle.Height(h - textAreaStyle.GetVerticalBorderSize())
		m.textArea.SetWidth(w - textAreaStyle.GetHorizontalBorderSize())
		m.textArea.SetHeight(h - textAreaStyle.GetHorizontalBorderSize())
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Submit):
			cmds = append(cmds, m.sendRequest())
		case key.Matches(msg, m.keys.ClearPrompt):
			m.textArea.SetValue("")
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}

	case respMsg:
		v := m.textArea.Value()
		v += "\n\n" + string(msg)
		m.textArea.SetValue(v)

	case errMsg:
		m.textArea.SetValue(msg.Error())
	}

	m.textArea, cmd = m.textArea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	b := strings.Builder{}

	title := titleStyle.Render(m.title) + "\n"
	textarea := textAreaStyle.Render(m.textArea.View())
	helpstr := helpStyle.Render(strings.Join(m.help, "\n"))
	versionInfo := versionInfoStyle.Render("Version: " + Version + "\n" + "Commit: " + Commit)

	b.WriteString(title)

	helpAndVersion := lipgloss.JoinVertical(lipgloss.Left, helpstr, versionInfo)
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, textarea, helpAndVersion))

	return b.String()
}

func (m Model) sendRequest() tea.Cmd {
	req := &ggpt3.CompletionRequest{
		Model:     ggpt3.TextDavinci003,
		Prompt:    m.textArea.Value(),
		MaxTokens: 500,
		TopP:      1,
	}

	resp, err := m.client.RequestCompletion(context.Background(), req)
	if err != nil {
		return func() tea.Msg {
			return errMsg(err)
		}
	}

	return func() tea.Msg {
		return respMsg(strings.Trim(resp.Choices[0].Text, "\n"))
	}
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
