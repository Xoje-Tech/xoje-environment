package styles

import "github.com/charmbracelet/lipgloss"

// Rose Pine Palette
var (
	Base    = lipgloss.Color("#191724")
	Surface = lipgloss.Color("#1f1d2e")
	Text    = lipgloss.Color("#e0def4")
	Love    = lipgloss.Color("#eb6f92") // Error
	Gold    = lipgloss.Color("#f6c177") // Warn
	Pine    = lipgloss.Color("#31748f") // Success/Info
	Foam    = lipgloss.Color("#9ccfd8")
	Iris    = lipgloss.Color("#c4a7e7") // Title
	Muted   = lipgloss.Color("#6e6a86")
)

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(Iris).
			Bold(true).
			MarginBottom(1)

	SuccessStyle = lipgloss.NewStyle().Foreground(Pine)
	WarnStyle    = lipgloss.NewStyle().Foreground(Gold)
	ErrorStyle   = lipgloss.NewStyle().Foreground(Love)
	MutedStyle   = lipgloss.NewStyle().Foreground(Muted)
	IrisStyle    = lipgloss.NewStyle().Foreground(Iris)
	GoldStyle    = lipgloss.NewStyle().Foreground(Gold)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Iris).
			Padding(1, 2).
			MarginTop(1)
)
