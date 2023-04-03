package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var title_margin_left = 1
var title_margin_top = 1
var title_margin_bottom = 1

var title_style_normal = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("0")).
  Background(lipgloss.Color("4")).
  MarginTop(title_margin_top).
  MarginBottom(title_margin_bottom).
  MarginLeft(title_margin_left)

var title_style_new_item = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("0")).
  Background(lipgloss.Color("5")).
  MarginTop(title_margin_top).
  MarginBottom(title_margin_bottom).
  MarginLeft(title_margin_left)

var title_style_delete = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("0")).
  Background(lipgloss.Color("1")).
  MarginTop(title_margin_top).
  MarginBottom(title_margin_bottom).
  MarginLeft(title_margin_left)

var title_style_edit = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("0")).
  Background(lipgloss.Color("3")).
  MarginTop(title_margin_top).
  MarginBottom(title_margin_bottom).
  MarginLeft(title_margin_left)

var normal_style = lipgloss.NewStyle().
  Bold(false).
  Foreground(lipgloss.Color("7"))

var selected_style = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("15"))

var check_style = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("2"))

var checked_style = lipgloss.NewStyle().
  Bold(false).
  Foreground(lipgloss.Color("#666"))

var delete_style = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("1"))

var edit_text_style = lipgloss.NewStyle().
  Bold(false).
  Background(lipgloss.Color("3")).
  Foreground(lipgloss.Color("1"))

var help_style_normal = lipgloss.NewStyle().
  Foreground(lipgloss.Color("#555"))

var help_style_key = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("#888"))
