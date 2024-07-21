package main

import (
	"fmt"
	"net"
	"log"
 "strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const WIDTH = 60
const HEIGHT = 5

type (
	errMsg error
)

// model stores application state
type model struct {
	conn *net.Conn
	viewport viewport.Model
	messages []string
	messageToSend string
	connected bool
	textarea textarea.Model
	senderStyle lipgloss.Style
	err error
}

type KeyMap struct {
	PageUp key.Binding
	PageDown key.Binding
}

// initialize
func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = WIDTH * HEIGHT 

	ta.SetWidth(WIDTH)
	ta.SetHeight(1)

	// remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(WIDTH, HEIGHT)
	vp.SetContent(`Welcome to TeaChat! Type a message and press enter to send.`)
	vp.KeyMap = viewport.KeyMap{
		PageUp: key.NewBinding( key.WithKeys("up"), key.WithHelp("↑", "scroll up message viewport")), 
		PageDown: key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "scroll down message viewport")),
	}

	return model{
		conn: nil,
		textarea: ta,
		messages: []string{},
		messageToSend: "",
		connected: false,
		viewport: vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err: nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textarea.Value() == "/connect" {
				m.conn = ConnectToServer()
				KeepAlive(*m.conn)
				m.connected = true
			}
			m.messageToSend = m.textarea.Value() 
			m.messages = append(m.messages, m.ChopText(m.senderStyle.Render("You: ") + m.messageToSend, WIDTH)...)
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()


			if m.connected {
				return m, SendCmd(m, *m.conn) 
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd) 
}

func (m model) View() string {
	return fmt.Sprintf("%s\n\n%s\n\n", m.viewport.View(), m.textarea.View(),)
}


func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("An error occured when attempting to run program: %v", err)
	}
}

/*
TODO:
	-[x] chop! the string array into [MAX_CHAR_WIDTH] slices (append the username first, then chop)
	-[x] don't let hjkl etc move the viewport 
	-[x] stay connected to server for program lifetime
	-[x] send message to server (simple)
	-[ ] send message struct to server (complex) using glob
		*/
