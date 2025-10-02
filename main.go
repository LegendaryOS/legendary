package main

import (
	"fmt"
	"os"
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

	inputStyle = lipgloss.NewStyle().
		Foreground(brightWhite).
		Background(darkGray).
		Padding(1, 2).
		MarginTop(1).
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(neonCyan).
		Width(70)
)

// Model for BubbleTea
type model struct {
	choices         []string
	cursor          int
	selected        bool
	enteringPackage bool
	packageName     string
	output          string
	running         bool
	errorMsg        string
	commandStatus   string
	tickCount       int // For animation
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
		if m.enteringPackage {
			switch msg.Type {
			case tea.KeyEnter:
				if m.packageName != "" || m.choices[m.cursor] == "update" {
					m.enteringPackage = false
					m.running = true
					m.commandStatus = "Executing..."
					return m, m.runCommand(m.choices[m.cursor], m.packageName)
				}
			case tea.KeyEsc:
				m.enteringPackage = false
				m.packageName = ""
				m.selected = false
			case tea.KeyBackspace:
				if len(m.packageName) > 0 {
					m.packageName = m.packageName[:len(m.packageName)-1]
				}
			case tea.KeyRunes:
				m.packageName += string(msg.Runes)
			}
			return m, nil
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
				if m.choices[m.cursor] == "install" || m.choices[m.cursor] == "remove" {
					m.enteringPackage = true
					m.packageName = ""
				} else {
					m.running = true
					m.commandStatus = "Executing..."
					return m, m.runCommand(m.choices[m.cursor], "")
				}
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
	// Header with simple styling
	header := lipgloss.NewStyle().
		Foreground(neonCyan).
		Bold(true).
		SetString("LegendaryOS").
		Render() + "\n" +
		subHeaderStyle.Render("Legendary Package Manager")

	// Animated border effect for header
	borderColor := neonPurple
	if m.tickCount%2 == 0 {
		borderColor = neonPink
	}
	animatedHeaderStyle := headerStyle.Copy().BorderForeground(borderColor)
	s := animatedHeaderStyle.Render(header)

	// Menu, package input, or output
	if !m.selected {
		s += "\n\nSelect Operation:\n\n"
		for i, choice := range m.choices {
			cursor := " "
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
	} else if m.enteringPackage {
		s += "\n\nOperation: " + selectedItemStyle.Render(strings.ToUpper(m.choices[m.cursor])) + "\n"
		s += "\nEnter package name:\n"
		s += inputStyle.Render(m.packageName + "_")
	} else {
		s += "\n\nOperation: " + selectedItemStyle.Render(strings.ToUpper(m.choices[m.cursor])) + "\n"
		if m.choices[m.cursor] != "update" {
			s += "\nPackage: " + inputStyle.Render(m.packageName) + "\n"
		}
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
	} else if m.enteringPackage {
		instructions += instructionStyle.Render("Type package name | Confirm: Enter | Cancel: Esc | Quit: q")
	} else {
		instructions += instructionStyle.Render("Return: Esc | Quit: q")
	}
	s += instructions

	return lipgloss.NewStyle().Margin(2, 4).Render(s) + "\n"
}

// Run the selected command
func (m model) runCommand(command, packageName string) tea.Cmd {
	return func() tea.Msg {
		switch command {
		case "install":
			if packageName == "" {
				return commandResult{err: fmt.Errorf("package name required for install")}
			}
			// Simulate command execution for demo (remove this in production)
			time.Sleep(2 * time.Second)
			output := fmt.Sprintf("Installing package %s...\nSample output for install command.", packageName)
			// Uncomment below for actual command execution
			/*
				cmd := exec.Command("/usr/lib/legendaryos/rpm-ostree", "install", packageName)
				out, err := cmd.CombinedOutput()
				if err != nil {
					return commandResult{err: err}
				}
				output := string(out)
			*/
			return commandResult{output: output, err: nil}
		case "remove":
			if packageName == "" {
				return commandResult{err: fmt.Errorf("package name required for remove")}
			}
			// Simulate command execution for demo (remove this in production)
			time.Sleep(2 * time.Second)
			output := fmt.Sprintf("Removing package %s...\nSample output for remove command.", packageName)
			// Uncomment below for actual command execution
			/*
				cmd := exec.Command("/usr/lib/legendaryos/rpm-ostree", "remove", packageName)
				out, err := cmd.CombinedOutput()
				if err != nil {
					return commandResult{err: err}
				}
				output := string(out)
			*/
			return commandResult{output: output, err: nil}
		case "update":
			// Simulate command execution for demo (remove this in production)
			time.Sleep(2 * time.Second)
			output := "Processing update operation...\nSample output for update command."
			// Uncomment below for actual command execution
			/*
				cmd := exec.Command("/bin/sh", "-c", "/usr/lib/legendaryos/rpm-ostree update && /usr/lib/legendaryos/rpm-ostree upgrade")
				out, err := cmd.CombinedOutput()
				if err != nil {
					return commandResult{err: err}
				}
				output := string(out)
			*/
			return commandResult{output: output, err: nil}
		default:
			return commandResult{err: fmt.Errorf("unknown command: %s", command)}
		}
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
