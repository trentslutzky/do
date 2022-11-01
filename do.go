package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	//"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
  "github.com/charmbracelet/bubbles/textinput"
)

type TodoListItem struct {
  Id int
  Title string
  Done bool
}

type model struct {
  items              []TodoListItem
  cursor             int
  selected           map[int]struct{}
  deleting           bool
  item_to_delete     int
  editing            bool
  creating           bool
  show_help          bool
  new_item_textInput textinput.Model
}

type TodoListItems struct {
  Items []TodoListItem
}

var plus_style = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("5"))

const (
  plus_icon = " "
  uncheck_icon = ""
  check_icon = ""
  edit_icon = "ﯽ"
  delete_icon = ""
)

func initialModel() model {

  var home = os.Getenv("HOME")
  var filename = home+"/.local/share/do/items.json"
  file, _ := ioutil.ReadFile(filename)
  data := TodoListItems{}
  fmt.Printf("%v",data)
  _ = json.Unmarshal([]byte(file),&data)

  ti := textinput.New()
  ti.Placeholder = "New todolist item"
  ti.Prompt = plus_style.Render(" "+plus_icon+" ")

  return model{
    items: data.Items,
    cursor: 0,
    deleting: false,
    editing: false,
    creating: false,
    show_help: true,
    new_item_textInput: ti,
  }
}

func (m model) Init() tea.Cmd {
  return textinput.Blink
}

func insert(a []TodoListItem, index int, value TodoListItem) []TodoListItem {
    if len(a) == index { // nil or empty slice or after last element
        return append(a, value)
    }
    a = append(a[:index+1], a[index:]...) // index < len(a)
    a[index] = value
    return a
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

  var is_not_editing = (m.creating == false && m.editing == false && m.deleting == false)

  var cmd tea.Cmd 
  switch msg := msg.(type) {

  case tea.KeyMsg:

    switch msg.String() {

    case "ctrl+c":
      return m, tea.Quit

    case "q":
      if is_not_editing {
        return m, tea.Quit
      }

    case "up":
      if m.cursor > 0 && is_not_editing {
        m.cursor--
      }

    case "down":
      if m.cursor < len(m.items)-1 && is_not_editing {
        m.cursor++
      }

    case "n":
      if is_not_editing {
        m.creating = true
      }

    case "e":
      if is_not_editing {
        m.editing = true
      }

    case "d":
      if m.deleting == true && len(m.items) > 0 {
        m.items = append(m.items[:m.cursor], m.items[m.cursor+1:]...)
        if m.cursor > 0 {
          m.cursor--
        }
        m.deleting = false
      }
      if is_not_editing && len(m.items) > 0 {
        m.deleting = true
        m.item_to_delete = m.cursor
      }

    case "?":
      if is_not_editing {
        m.show_help = !m.show_help
      }

    case "esc":
      if is_not_editing == false {
        m.creating = false
        m.editing = false
        m.deleting = false
        m.new_item_textInput.Blur()
        _ = m.new_item_textInput.Reset()
      }

    case " ":
      if is_not_editing {
        m.items[m.cursor].Done = !m.items[m.cursor].Done
      }

    case "enter":
      if m.creating && len(m.new_item_textInput.Value()) > 0 {
        new_item := TodoListItem{
          Title: m.new_item_textInput.Value(),
          Done: false,
        }
        if len(m.items) > 0 {
          m.items = insert(m.items, m.cursor + 1, new_item)
          m.cursor++
        } else {
          m.items = append(m.items, new_item)
        }
        m.new_item_textInput.Blur()
        _ = m.new_item_textInput.Reset()
        m.creating = false
      } else {
        if is_not_editing && len(m.items) > 0 {
          m.items[m.cursor].Done = !m.items[m.cursor].Done
        }
      }
    }
  }

  if m.creating == true {
    if m.new_item_textInput.Focused() == true {
      m.new_item_textInput, cmd = m.new_item_textInput.Update(msg)
    } else {
      m.new_item_textInput.Focus()
    }
  }

  return m, cmd
}

var title_style_normal = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("0")).
  Background(lipgloss.Color("4")).
  MarginTop(0).
  MarginBottom(1).
  MarginLeft(0)

var title_style_new_item = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("0")).
  Background(lipgloss.Color("5")).
  MarginTop(0).
  MarginBottom(1).
  MarginLeft(0)

var title_style_delete = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("0")).
  Background(lipgloss.Color("1")).
  MarginTop(0).
  MarginBottom(1).
  MarginLeft(0)

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

var help_style_normal = lipgloss.NewStyle().
  Foreground(lipgloss.Color("#555"))

var help_style_key = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("#888"))

func (m model) View() string {

  const header_text = "             Todo List              "

  s := ""
  if m.creating {
    s += title_style_new_item.Render(header_text)
  } else if m.deleting {
    s += title_style_delete.Render(header_text)
  } else {
    s += title_style_normal.Render(header_text)
  }
  s += "\n"

  for i, item := range m.items {

    line_string := ""
    checked := " "

    if item.Done {
      checked = " " + check_style.Render(check_icon) + "  "
      if m.cursor == i && m.creating == false {
        line_string = selected_style.Render("%s")
      } else {
        line_string = checked_style.Render("%s")
      }
    } else {
      if m.cursor == i && m.creating == false {
        checked = " " + uncheck_icon + "  "
        line_string = selected_style.Render("%s")
      } else {
        checked = checked_style.Render(" " + uncheck_icon + "  ")
        line_string = normal_style.Render("%s")
      }
    }
    if m.deleting && m.item_to_delete == i {
      checked = delete_style.Render(" " + delete_icon + "  ")
      line_string = delete_style.Render("%s")
    }

    s += fmt.Sprintf(checked)
    s += fmt.Sprintf(line_string, item.Title)
    if m.cursor == i && m.creating {
      s += "\n"
      s += fmt.Sprintf("%s", m.new_item_textInput.View())
    }
    s += "\n"
  }
  
  if len(m.items) == 0 {
    if m.creating {
      s += fmt.Sprintf("%s", m.new_item_textInput.View())
      s += "\n\n"
    }
  } else {
    s += "\n"
  }

  if m.show_help {
    if m.deleting {
      s += " "
      s += help_style_key.Render("d ")
      s += help_style_normal.Render("confirm")
      s += "  "
      s += help_style_key.Render("esc ")
      s += help_style_normal.Render("cancel")
      s += "  "
    } else if m.creating || m.editing {
      s += " "
      s += help_style_key.Render("esc ")
      s += help_style_normal.Render("cancel")
      s += "  "
      s += help_style_key.Render(" ")
      s += help_style_normal.Render("save")
      s += "  "
    } else {
      s += " "
      s += help_style_key.Render("↑/↓ ")
      s += help_style_normal.Render("move")
      s += "  "
      s += help_style_key.Render("q ")
      s += help_style_normal.Render("quit")
      s += "  "
      s += help_style_key.Render("n ")
      s += help_style_normal.Render("new")
      s += "  "
      s += help_style_key.Render("d ")
      s += help_style_normal.Render("delete")
      s += "\n "
      s += help_style_key.Render(" ")
      s += help_style_normal.Render("toggle")
      s += "  "
      s += help_style_key.Render("? ")
      s += help_style_normal.Render("toggle help")
    }
  }

  return s
}

func main() {
  p := tea.NewProgram(initialModel(), tea.WithAltScreen())
  if err := p.Start(); err != nil {
    fmt.Printf("error %v", err)
    os.Exit(1)
  }
}
