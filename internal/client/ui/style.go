package ui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

// Цвета
var (
	colorWhite = lipgloss.Color("#EEEEEE")
	colorGray  = lipgloss.Color("#575757")
	colorPink  = lipgloss.Color("#d94cc6")
	colorGreen = lipgloss.Color("#008000")
	colorRed   = lipgloss.Color("#FF0000")
)

// Стили
var (
	headerStyle  = lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Foreground(colorWhite).Background(colorGray).Bold(true)
	inputStyle   = lipgloss.NewStyle().Foreground(colorPink).Bold(true).Blink(true)
	outputStyle  = lipgloss.NewStyle().Width(150).Align(lipgloss.Left).Bold(true).Foreground(colorGreen)
	successStyle = lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Foreground(colorWhite).Background(colorGreen)
	errorStyle   = lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Foreground(colorWhite).Background(colorRed).Bold(true)
)

// Стили для таблицы
var (
	tHeaderStyle  = lipgloss.NewRenderer(os.Stdout).NewStyle().Foreground(colorPink).Bold(true).Align(lipgloss.Center)
	tCellStyle    = lipgloss.NewRenderer(os.Stdout).NewStyle().Padding(0, 1).Width(14).Align(lipgloss.Center)
	tOddRowStyle  = tCellStyle.Foreground(colorGray)
	tEvenRowStyle = tCellStyle.Foreground(colorWhite)
	tBorderStyle  = lipgloss.NewStyle().Foreground(colorPink)
)
