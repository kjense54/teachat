package main

import (
	"fmt"
	"context"
	"net"
	"log"
	"time"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const WIDTH = 30
const HEIGHT = 5

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
		textarea: ta,
		messages: []string{},
		viewport: vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err: nil,
	}
}

func (m model) Init() tea.Cmd {
	//connectToServer() 
	return textarea.Blink
}

func (m model) chopText(text string, size int) []string {
	if len(text) == 0 {
		return nil
	}
	if len(text) < size {
		return []string{text}
	}
	var chopped []string = make([]string, 0, (len(text)-1)/size+1)
	currentLen := 0
	currentStart := 0
	for i := range text {
		if currentLen == size {
			chopped = append(chopped, text[currentStart:i])
			currentLen = 0
			currentStart = i
		} 
		currentLen++
	}
	// add extra bits at end
	chopped = append(chopped, text[currentStart:])
	return chopped
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
			concat := m.senderStyle.Render("You: ") + m.textarea.Value() 
			m.messages = append(m.messages, m.chopText(concat, WIDTH)...)
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()

		/*
TODO:
			-[ ] read/write to buffer of specified size to ensure message recieved correctly
			-[ ] chop! the string array into [MAX_CHAR_WIDTH] slices (append the username first, then chop)
			-[ ] don't let hjkl etc move the viewport 
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
	return fmt.Sprintf("%s\n\n%s\n\n", m.viewport.View(), m.textarea.View(),)
}


func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("An error occured when attempting to run program: %v", err)
	}
}
