package styles

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"reflect"
)

// TODO - standardize the widths and heights of everything!

const (
	// BaseWidth represents the basic size of a "single" BaseWidth
	BaseWidth = 60
	maxWidth  = 120
	// BaseHeight represents the basic size of a "single" BaseHeight
	BaseHeight = 16
	maxHeight  = 40
)

const (
	Wide = iota
	Tall
	Single
	Full
)

type Size = int

type BoxSize struct {
	Width  int
	Height int
}

type BoxOptions struct {
	// Specify a preset Size option
	Size Size
	// Specify the Box BorderColor
	BorderColor lipgloss.Color
	TextColor   lipgloss.Color
	// Specify the exact BoxSize
	BoxSize BoxSize
}

var (
	baseBoxStyle = lipgloss.NewStyle().Padding(2).Border(lipgloss.RoundedBorder())
)

func Box(options BoxOptions) lipgloss.Style {
	if reflect.ValueOf(options.BorderColor).IsZero() {
		options.BorderColor = SecondaryGray
	}

	if !reflect.ValueOf(options.BoxSize).IsZero() {
		return baseBoxStyle.Copy().Width(options.BoxSize.Width).Height(options.BoxSize.Height).BorderForeground(options.BorderColor)
	}

	var style lipgloss.Style
	switch options.Size {
	case Wide:
		return baseBoxStyle.Copy().Width(BaseWidth * 2).Height(BaseHeight).BorderForeground(options.BorderColor)
	case Tall:
		return baseBoxStyle.Copy().Width(BaseWidth).Height(BaseHeight * 2).BorderForeground(options.BorderColor)
	case Single:
		return baseBoxStyle.Copy().Width(BaseWidth).Height(BaseHeight).BorderForeground(options.BorderColor)
	case Full:
		return baseBoxStyle.Copy().Width(BaseWidth * 2).Height(BaseHeight * 2).BorderForeground(options.BorderColor)
	}

	if !reflect.ValueOf(options.TextColor).IsZero() {
		return style.Foreground(options.TextColor)
	}

	return style
}

var (
	//WindowStyle    = lipgloss.NewStyle().Padding(2).Width(100).Height(45).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b"))
	CenterStyle    = lipgloss.NewStyle().Align(lipgloss.Center)
	FormTitleStyle = lipgloss.NewStyle().Foreground(SecondaryColor)
	FormLabelStyle = lipgloss.NewStyle().Foreground(SecondaryGray)
	FormErrorStyle = lipgloss.NewStyle().Foreground(PrimaryColor)
	FocusedStyle   = lipgloss.NewStyle().Foreground(AccentColor)

	WarnToastStyle = lipgloss.NewStyle().Foreground(PrimaryGray).Border(lipgloss.RoundedBorder()).BorderForeground(PrimaryColor).Width(25).Align(lipgloss.Center)
	InfoToastStyle = lipgloss.NewStyle().Foreground(PrimaryGray).Border(lipgloss.RoundedBorder()).BorderForeground(SecondaryColor).Width(25).Align(lipgloss.Center)

	SelectedItemStyle   = lipgloss.NewStyle().Foreground(PrimaryColor).Width(60)
	DefaultItemStyle    = lipgloss.NewStyle().Foreground(SecondaryGray).Width(60)
	ActiveCursorStyle   = lipgloss.NewStyle().Foreground(PrimaryColor).Width(2)
	InactiveCursorStyle = lipgloss.NewStyle().Foreground(SecondaryGray).Width(2)

	PaginationStyle = list.DefaultStyles().PaginationStyle

	InactiveStyle = lipgloss.NewStyle().Foreground(SecondaryGray)

	PageTitleStyle       = lipgloss.NewStyle().Foreground(SecondaryGray).Bold(true)
	ActivePageTitleStyle = lipgloss.NewStyle().Foreground(SecondaryColor).Bold(true)
)
