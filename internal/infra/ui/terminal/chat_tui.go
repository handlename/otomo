package terminal

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
)

type thinkResultMsg struct {
	ans *reasoning.Answer
	err error
}

type model struct {
	ctx          context.Context
	otomo        *chat.Otomo
	tools        []reasoning.Tool
	history      []*core.Message
	viewport     viewport.Model
	textInput    textinput.Model
	spinner      spinner.Model
	thinking     bool
	userInputVal string
	err          error
}

func StartChatTUI(ctx context.Context, otomo *chat.Otomo, tools []reasoning.Tool) error {
	ti := textinput.New()
	ti.Placeholder = "Ask Otomo..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	vp := viewport.New(80, 20)
	vp.SetContent("Welcome to Otomo Terminal Chat!\nPress Enter to send. Ctrl+C to quit.\n")

	m := model{
		ctx:       ctx,
		otomo:     otomo,
		tools:     tools,
		viewport:  vp,
		textInput: ti,
		spinner:   s,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			if m.thinking {
				return m, nil
			}
			input := m.textInput.Value()
			if strings.TrimSpace(input) == "" {
				return m, nil // Skip empty input validation (ignores empty enter keys)
			}

			m.userInputVal = input
			m.textInput.SetValue("")
			m.thinking = true

			// Append user prompt to viewport
			m.appendToHistory(core.RoleUser, input)
			m.textInput.Blur()

			return m, tea.Batch(
				m.spinner.Tick,
				m.thinkCmd(input),
			)
		}

	case tea.WindowSizeMsg:
		// Make it responsive to window resize
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4 // Leave room for prompt input, spinner, and exit instructions
		if m.viewport.Height < 0 {
			m.viewport.Height = 0
		}

	case spinner.TickMsg:
		if m.thinking {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case thinkResultMsg:
		m.thinking = false
		m.textInput.Focus()
		if msg.err != nil {
			m.appendToHistory(core.RoleAssistant, fmt.Sprintf("Error: %v", msg.err))
		} else {
			m.appendToHistory(core.RoleAssistant, string(msg.ans.Body()))
			// Add to structural history
			uMsg, _ := core.NewMessage(core.RoleUser, core.UserID{}, core.MessageBody(m.userInputVal))
			aiMsg, _ := core.NewMessage(core.RoleAssistant, core.UserID{}, core.MessageBody(msg.ans.Body()))
			m.history = append(m.history, uMsg, aiMsg)
		}
	}

	if !m.thinking {
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) appendToHistory(role core.MessageRole, content string) {
	style := lipgloss.NewStyle()
	prefix := ""
	if role == core.RoleUser {
		style = style.Foreground(lipgloss.Color("2")).Bold(true)
		prefix = "You: "
	} else {
		style = style.Foreground(lipgloss.Color("4")).Bold(true)
		prefix = "Otomo: "
	}

	current := m.viewport.View()
	lines := strings.Split(current, "\n")
	lines = append(lines, style.Render(prefix)+content, "")
	m.viewport.SetContent(strings.Join(lines, "\n"))
	m.viewport.GotoBottom()
}

func (m model) thinkCmd(prompt string) tea.Cmd {
	return func() tea.Msg {
		c := reasoning.NewContext()
		if len(m.history) > 0 {
			_ = c.SetMessages(m.history)
		}
		c.SetUserPrompt(core.PromptBody(prompt))
		c.SetTools(m.tools)

		ans, err := executeToolLoop(m.ctx, m.otomo, c, m.tools)
		return thinkResultMsg{ans: ans, err: err}
	}
}

func (m model) View() string {
	var s strings.Builder
	s.WriteString(m.viewport.View())
	s.WriteString("\n")

	if m.thinking {
		s.WriteString(m.spinner.View() + " Otomo is thinking...\n")
	} else {
		s.WriteString(m.textInput.View())
		s.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("Enter to send • Esc to exit"))
	}
	return s.String()
}
