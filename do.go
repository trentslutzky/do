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
  item_to_edit       int
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

var edit_style = lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("3"))

const (
  plus_icon = " "
  uncheck_icon = ""
  check_icon = ""
  edit_icon = " "
  delete_icon = ""
)

func initialModel() model {

  var home = os.Getenv("HOME")
  var filename = home+"/.local/share/do/items.json"
  file, _ := ioutil.ReadFile(filename)
  data := TodoListItems{}
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
    item_to_edit: -1,
    item_to_delete: -1,
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
        var home = os.Getenv("HOME")
        var filename = home+"/.local/share/do/items.json"
        new_items := TodoListItems{
          Items: m.items,
        }
        file, _ := json.MarshalIndent(new_items, "", " ")
        _ = ioutil.WriteFile(filename,file, 0644)
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

    case "n", "o":
      if is_not_editing {
        m.new_item_textInput.Prompt = plus_style.Render(" "+plus_icon+" ")
        m.creating = true
      }

    case "e":
      if is_not_editing && len(m.items) > 0 {
        m.editing = true
        m.new_item_textInput.Prompt = edit_style.Render(" "+edit_icon+" ")
        m.new_item_textInput.SetValue(m.items[m.cursor].Title)
      }

    case "G":
      if is_not_editing && len(m.items) > 0 {
        m.cursor = len(m.items) - 1
      }

    case "g":
      if is_not_editing && len(m.items) > 0 {
        m.cursor = 0
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
        m.item_to_delete = -1
        m.item_to_edit = -1
        m.new_item_textInput.Blur()
        _ = m.new_item_textInput.Reset()
      }

    case " ":
      if is_not_editing {
        m.items[m.cursor].Done = !m.items[m.cursor].Done
      }

    case "enter":
      if m.editing {
        m.items[m.cursor].Title = m.new_item_textInput.Value()
        m.new_item_textInput.Blur()
        _ = m.new_item_textInput.Reset()
        m.editing = false
        m.item_to_edit = -1
      } else if m.creating {
        if len(m.new_item_textInput.Value()) > 0 {
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
        }
      } else {
        if is_not_editing && len(m.items) > 0 {
          m.items[m.cursor].Done = !m.items[m.cursor].Done
        }
      }
    }
  }

  if m.editing && m.item_to_edit == -1 {
    m.item_to_edit = m.cursor
    m.new_item_textInput.Update(m.items[m.cursor].Title)
  }

  if m.creating == true || m.editing == true {
    if m.new_item_textInput.Focused() == true {
      m.new_item_textInput, cmd = m.new_item_textInput.Update(msg)
    } else {
      m.new_item_textInput.Focus()
    }
  }

  return m, cmd
}

var title_margin_left = 2
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

func (m model) View() string {

  const header_text = "             Todo List              "

  s := ""
  if m.creating {
    s += title_style_new_item.Render(header_text)
  } else if m.deleting {
    s += title_style_delete.Render(header_text)
  } else if m.editing {
    s += title_style_edit.Render(header_text)
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
    if m.editing && m.item_to_edit == i {
      checked = edit_style.Render(" " + delete_icon + "  ")
      line_string = edit_text_style.Render("%s")
    }

    if( m.item_to_edit != i ) {
      s += fmt.Sprintf(checked)
      s += fmt.Sprintf(line_string, item.Title)
      s += "\n"
    }
    if m.cursor == i && m.creating {
      s += fmt.Sprintf("%s", m.new_item_textInput.View())
      s += "\n"
    }
    if m.cursor == i && m.editing {
      m.new_item_textInput.Update(item.Title)
      s += fmt.Sprintf("%s", m.new_item_textInput.View())
      s += "\n"
    }
  }
  
  if len(m.items) == 0 {
    if m.creating {
      s += fmt.Sprintf("%s", m.new_item_textInput.View())
      s += "\n\n"
    }
  } else {
    s += "\n"
  }

  command := ""
  if m.editing {
    command = edit_style.Render(" edit")+" -"
  } else if m.deleting {
    command = delete_style.Render(" delete")+" -"
  } else if m.creating {
    command = plus_style.Render(" new")+" -"
  }

  s += command
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
      s += help_style_key.Render("n/o ")
      s += help_style_normal.Render("new")
      s += "  "
      s += help_style_key.Render("d ")
      s += help_style_normal.Render("delete")
      s += "\n "
      s += help_style_key.Render(" ")
      s += help_style_normal.Render("toggle")
      s += "  "
      s += help_style_key.Render("e ")
      s += help_style_normal.Render("edit")
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
