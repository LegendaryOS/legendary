package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Cyberpunk color palette
var (
	neonCyan    = lipgloss.Color("#00F7FF")
	neonPurple  = lipgloss.Color("#D600FF")
	neonPink    = lipgloss.Color("#FF007A")
	darkGray    = lipgloss.Color("#1A1A1A")
	brightWhite = lipgloss.Color("#E6E6E6")
	neonGreen   = lipgloss.Color("#00FF9F")
	neonBlue    = lipgloss.Color("#007BFF")
)

// Styles for the CLI
var (
	headerStyle = lipgloss.NewStyle().
			Foreground(neonCyan).
			Background(darkGray).
			Bold(true).
			Padding(1, 3).
			Border(lipgloss.DoubleBorder(), true).
			BorderForeground(neonPurple).
			Align(lipgloss.Center).
			Width(60).
			BorderBackground(neonBlue)

	subHeaderStyle = lipgloss.NewStyle().
			Foreground(neonGreen).
			Background(darkGray).
			Italic(true).
			Padding(0, 3).
			Width(60).
			Align(lipgloss.Center)

	itemStyle = lipgloss.NewStyle().
			Foreground(brightWhite).
			PaddingLeft(4).
			MarginTop(1)

	selectedItemStyle = lipgloss.NewStyle().
			Foreground(neonCyan).
			Background(neonPurple).
			PaddingLeft(4).
			MarginTop(1).
			Bold(true).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(neonPink)

	outputStyle = lipgloss.NewStyle().
			Foreground(brightWhite).
			Background(darkGray).
			Padding(1, 2).
			MarginTop(1).
			Border(lipgloss.ThickBorder(), true).
			BorderForeground(neonGreen).
			Width(70)

	errorStyle = lipgloss.NewStyle().
			Foreground(brightWhite).
			Background(neonPink).
			Padding(1, 2).
			MarginTop(1).
			Border(lipgloss.ThickBorder(), true).
			BorderForeground(neonCyan).
			Width(70)

	instructionStyle = lipgloss.NewStyle().
			Foreground(neonBlue).
			Padding(1, 2).
			Align(lipgloss.Left)

	statusStyle = lipgloss.NewStyle().
			Foreground(neonGreen).
			Background(darkGray).
			Padding(0, 1).
			MarginTop(1)

	// Gradient effect for header
	gradientStyle = lipgloss.NewStyle().
			Foreground(neonCyan).
			SetString("LegendaryOS").
			Faint(false).
			Gradient(lipgloss.GradientColors{
				{Color: neonCyan, Offset: 0},
				{Color: neonPurple, Offset: 0.5},
				{Color: neonPink, Offset: 1},
			})
)

// Model for BubbleTea
type model struct {
	choices       []string
	cursor        int
	selected      bool
	output        string
	running       bool
	errorMsg      string
	commandStatus string
	tickCount     int // For animation
}

// Initialize the model
func initialModel() model {
	return model{
		choices:  []string{"install", "remove", "update"},
		cursor:   0,
		selected: false,
	}
}

// BubbleTea Init function
func (m model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

// Tick message for animations
type tickMsg struct{}

// Command result message
type commandResult struct {
	output string
	err    error
}

// BubbleTea Update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.running {
			return m, nil // Ignore input while command is running
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if !m.selected {
				m.selected = true
				m.running = true
				m.commandStatus = "Executing..."
				return m, m.runCommand(m.choices[m.cursor])
			}
		case "esc":
			if m.selected {
				m.selected = false
				m.output = ""
				m.errorMsg = ""
				m.commandStatus = ""
			}
		}

	case commandResult:
		m.running = false
		m.commandStatus = "Execution Complete"
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Error: %v", msg.err)
			m.output = ""
		} else {
			m.output = msg.output
			m.errorMsg = ""
		}

	case tickMsg:
		m.tickCount++
		return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	}

	return m, nil
}

// BubbleTea View function
func (m model) View() string {
	// Header with gradient effect
	header := gradientStyle.Render("LegendaryOS") + "\n" +
		subHeaderStyle.Render("Legendary Package Manager")

	// Animated border effect for header
	borderColor := neonPurple
	if m.tickCount%2 == 0 {
		borderColor = neonPink
	}
	animatedHeaderStyle := headerStyle.Copy().BorderForeground(borderColor)
	s := animatedHeaderStyle.Render(header)

	// Menu or output
	if !m.selected {
		s += "\n\nSelect Operation:\n\n"
		for i, choice := range m.choices {
			cursor := "  "
			if m.cursor == i {
				cursor = ">> "
			}
			if m.cursor == i {
				s += selectedItemStyle.Render(cursor + strings.ToUpper(choice))
			} else {
				s += itemStyle.Render(cursor + strings.ToUpper(choice))
			}
			s += "\n"
		}
	} else {
		s += "\n\nOperation: " + selectedItemStyle.Render(strings.ToUpper(m.choices[m.cursor])) + "\n"
		s += "\nStatus: " + statusStyle.Render(m.commandStatus) + "\n"
		if m.errorMsg != "" {
			s += "\n" + errorStyle.Render(m.errorMsg)
		} else if m.output != "" {
			s += "\n" + outputStyle.Render(m.output)
		}
	}

	// Instructions
	instructions := "\n"
	if !m.selected {
		instructions += instructionStyle.Render("Navigate: ↑↓ | Select: Enter | Quit: q")
	} else {
		instructions += instructionStyle.Render("Return: Esc | Quit: q")
	}
	s += instructions

	return lipgloss.NewStyle().Margin(2, 4).Render(s) + "\n"
}

// Run the selected command
func (m model) runCommand(command string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		switch command {
		case "install":
			cmd = exec.Command("/usr/lib/legendaryos/rpm-ostree", "install")
		case "remove":
			cmd = exec.Command("/usr/lib/legendaryos/rpm-ostree", "remove")
		case "update":
			cmd = exec.Command("/bin/sh", "-c", "/usr/lib/legendaryos/rpm-ostree update && /usr/lib/legendaryos/rpm-ostree upgrade")
		default:
			return commandResult{err: fmt.Errorf("unknown command: %s", command)}
		}

		// Simulate command execution for demo (remove this in production)
		time.Sleep(2 * time.Second)
		output := fmt.Sprintf("Processing %s operation...\nSample output for %s command.", command, command)

		// Uncomment below for actual command execution
		/*
			out, err := cmd.CombinedOutput()
			if err != nil {
				return commandResult{err: err}
			}
			output := string(out)
		*/

		return commandResult{output: output, err: nil}
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
