package main

import (
	"fmt"
	"context"
	"net"
	"log"
	"time"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// dial the server
func connectToServer() net.Conn {
	address := "localhost:33183"
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", address)	
	if err != nil {
		log.Fatal("Failed to %v", err)
	}
	return conn
}

type (
	errMsg error
)

// model stores application state
type model struct {
	viewport viewport.Model
	messages []string
	textarea textarea.Model
	senderStyle lipgloss.Style
	err error
}

// initialize
func initialModel() model {
	const WIDTH = 30
	const HEIGHT = 5

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = 280

	ta.SetWidth(WIDTH)
	ta.SetHeight(1)

	// remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(WIDTH, HEIGHT)
	vp.SetContent(`Type a message and press Enter to chat`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea: ta,
		messages: []string{},
		viewport: vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err: nil,
	}
}

func (m model) Init() tea.Cmd {
	connectToServer() 
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
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		// TODO: read from buffer more precisely, input text not always showing up until enter is pressed multiple times
		/* TODO
			-[ ] read/write to buffer of specified size to ensure message recieved correctly
			-[ ] only allow input up to [MAX_CHAR_LIMIT]
			-[ ] send message to server (simple)
			-[ ] send message struct to server (complex) using glob
			*/
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd) 
}

func (m model) View() string {
	return fmt.Sprintf("%s\n\n%s", m.viewport.View(), m.textarea.View(),)+ "\n\n"
}


func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("An error occured when attempting to run program: %v", err)
	}
}
