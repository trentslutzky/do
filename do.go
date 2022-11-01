package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TodoListItem struct {
  Id int
  Title string
  Done bool
}

type model struct {
  items  []TodoListItem
  cursor   int
  selected map[int]struct{}
  deleting bool
}

type TodoListItems struct {
  Items []TodoListItem
}

func initialModel() model {
  var home = os.Getenv("HOME")
  var filename = home+"/.local/share/do/items.json"
  file, _ := ioutil.ReadFile(filename)
  data := TodoListItems{}
  _ = json.Unmarshal([]byte(file),&data)
  return model{
    items: data.Items,
    cursor: 0,
    deleting: false,
  }
}

func (m model) Init() tea.Cmd {
  return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {

  case tea.KeyMsg:

    switch msg.String() {

    case "ctrl+c", "q":
      return m, tea.Quit

    case "up":
      if m.cursor > 0 {
        m.cursor--
      }

    case "down":
      if m.cursor < len(m.items)-1 {
        m.cursor++
      }

    case "enter", " ":
      m.items[m.cursor].Done = !m.items[m.cursor].Done

    }
  }

  return m, nil
}

var title_style = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("0")).
  Background(lipgloss.Color("4")).
  MarginTop(1).
  MarginBottom(1).
  MarginLeft(2)

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

func (m model) View() string {
  s := title_style.Render(" Todo List ")
  s += "\n"
  sort.SliceStable(m.items, func(i, j int) bool {
    return m.items[i].Done == false
  })
  sort.SliceStable(m.items, func(i, j int) bool {
    return m.items[i].Id < m.items[j].Id
  })
  for i, item := range m.items {


    line_string := ""
    checked := " "

    if item.Done {
      checked = check_style.Render("   ")
      if m.cursor == i {
        line_string = selected_style.Render("%s")
      } else {
        line_string = checked_style.Render("%s")
      }
    } else {
      if m.cursor == i {
        checked = "   "
        line_string = selected_style.Render("%s")
      } else {
        checked = "   "
        line_string = normal_style.Render("%s")
      }
    }

    s += fmt.Sprintf(checked)
    s += fmt.Sprintf(line_string, item.Title)

    s += "\n"
  }
  s += "\n"

  return s
}

func main() {
  p := tea.NewProgram(initialModel(), tea.WithAltScreen())
  if err := p.Start(); err != nil {
    fmt.Printf("erroro %v", err)
    os.Exit(1)
  }
}
