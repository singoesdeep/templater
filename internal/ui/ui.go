package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// Colors for different message types
var (
	SuccessColor = color.New(color.FgGreen)
	ErrorColor   = color.New(color.FgRed)
	InfoColor    = color.New(color.FgBlue)
	WarnColor    = color.New(color.FgYellow)
)

// ProgressBar represents a progress indicator
type ProgressBar struct {
	total     int
	current   int
	width     int
	startTime time.Time
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		total:     total,
		width:     50,
		startTime: time.Now(),
	}
}

// Update updates the progress bar
func (p *ProgressBar) Update(current int) {
	p.current = current
	p.draw()
}

// Increment increments the progress bar
func (p *ProgressBar) Increment() {
	p.current++
	p.draw()
}

// draw draws the progress bar
func (p *ProgressBar) draw() {
	// Calculate progress
	progress := float64(p.current) / float64(p.total)
	width := int(float64(p.width) * progress)

	// Create progress bar
	bar := strings.Repeat("=", width) + strings.Repeat("-", p.width-width)

	// Calculate elapsed time
	elapsed := time.Since(p.startTime)

	// Print progress bar
	fmt.Printf("\r[%s] %d%% (%d/%d) %s", bar, int(progress*100), p.current, p.total, elapsed.Round(time.Second))
	if p.current == p.total {
		fmt.Println()
	}
}

// PrintSuccess prints a success message
func PrintSuccess(format string, args ...interface{}) {
	SuccessColor.Printf("✅ "+format+"\n", args...)
}

// PrintError prints an error message
func PrintError(format string, args ...interface{}) {
	ErrorColor.Printf("❌ "+format+"\n", args...)
}

// PrintInfo prints an info message
func PrintInfo(format string, args ...interface{}) {
	InfoColor.Printf("ℹ️  "+format+"\n", args...)
}

// PrintWarning prints a warning message
func PrintWarning(format string, args ...interface{}) {
	WarnColor.Printf("⚠️  "+format+"\n", args...)
}

// PromptYesNo prompts for yes/no input
func PromptYesNo(message string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return strings.ToLower(result) == "y", nil
}

// PromptSelect prompts for selection from a list
func PromptSelect(message string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: message,
		Items: items,
	}

	_, result, err := prompt.Run()
	return result, err
}

// PromptInput prompts for text input
func PromptInput(message string) (string, error) {
	prompt := promptui.Prompt{
		Label: message,
	}

	return prompt.Run()
}

// PrintTable prints a table with headers and rows
func PrintTable(headers []string, rows [][]string) {
	// Calculate column widths
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print headers
	for i, header := range headers {
		fmt.Printf("%-*s", widths[i]+2, header)
	}
	fmt.Println()

	// Print separator
	for _, width := range widths {
		fmt.Print(strings.Repeat("-", width+2))
	}
	fmt.Println()

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			fmt.Printf("%-*s", widths[i]+2, cell)
		}
		fmt.Println()
	}
}

// PrintDiff prints a diff between two strings
func PrintDiff(old, new string) {
	// Split into lines
	oldLines := strings.Split(old, "\n")
	newLines := strings.Split(new, "\n")

	// Print diff
	for i := 0; i < len(oldLines) || i < len(newLines); i++ {
		if i >= len(oldLines) {
			SuccessColor.Printf("+ %s\n", newLines[i])
		} else if i >= len(newLines) {
			ErrorColor.Printf("- %s\n", oldLines[i])
		} else if oldLines[i] != newLines[i] {
			ErrorColor.Printf("- %s\n", oldLines[i])
			SuccessColor.Printf("+ %s\n", newLines[i])
		} else {
			fmt.Printf("  %s\n", oldLines[i])
		}
	}
}

// PrintSpinner prints a loading spinner
func PrintSpinner(message string, done chan bool) {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Printf("\r%s\n", strings.Repeat(" ", len(message)+3))
			return
		default:
			fmt.Printf("\r%s %s", spinner[i], message)
			i = (i + 1) % len(spinner)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// PrintBox prints a message in a box
func PrintBox(message string) {
	lines := strings.Split(message, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	fmt.Println("┌" + strings.Repeat("─", maxLen+2) + "┐")
	for _, line := range lines {
		fmt.Printf("│ %-*s │\n", maxLen, line)
	}
	fmt.Println("└" + strings.Repeat("─", maxLen+2) + "┘")
}
