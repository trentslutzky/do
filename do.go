package main

import (
  "fmt"
  "os"
  "encoding/json"
  "io/ioutil"
  "strconv"
)

type TodoListItem struct {
  Title string
  Done bool
}

type TodoListItems struct {
  Items []TodoListItem
}

const filename = "/Users/trent/.local/share/do/items.json"

func main() {

  file, _ := ioutil.ReadFile(filename)
  data := TodoListItems{}

  _ = json.Unmarshal([]byte(file),&data)

  args := os.Args[1:]

  // new_item := TodoListItem{
  //   Title: "test",
  //   Done: false,
  // } 

  // data := TodoListItems{
  //   Items: []TodoListItem{new_item},
  // }

  // new_item_two := TodoListItem{
  //   Title: "test_two",
  //   Done: false,
  // } 

  // fmt.Println(data)

  //data.Items = data.Items.append(data.Items,new_item_two)

  // if no command is passed, print out the todo list.
  if len(args) == 0 {
    listItems(args,data)
    return
  }

  com := args[0]


  switch com {
  case "new":
    data = newItem(args,data)
  case "rename":
    data = renameItem(args,data)
  case "toggle":
    data = toggleItem(args,data)
  case "done":
    data = toggleItem(args,data)
  case "del":
    data = delItem(args,data)
  default:
    fmt.Println("Unrecognized command.")
  }

  file, _ = json.MarshalIndent(data, "", " ")
  _ = ioutil.WriteFile(filename,file, 0644)

}

func listItems(args []string, data TodoListItems) {
  fmt.Println("todo list")
  for i := 0; i < len(data.Items); i++{
    title := data.Items[i].Title
    toggleString := ""
    space := ""
    if len(data.Items) > 9 {
      if(i < 10) {
        space = " "
      }
    }
    if data.Items[i].Done {
      toggleString = "[x]"
    } else {
      toggleString = "[ ]"
    }
    fmt.Printf("%s%v %s %s\n",space,i,toggleString,title)
  }
}

func newItem(args []string, data TodoListItems) TodoListItems {
  title := args[1]
  new_item := TodoListItem{
    Title: title,
    Done: false,
  }
  data.Items = append(data.Items,new_item)
  listItems(args,data)
  return data
}

func delItem(args []string, data TodoListItems) TodoListItems {
  delIndex,_ := strconv.Atoi(args[1])
  data.Items = append(data.Items[:delIndex],data.Items[delIndex+1:]...)
  listItems(args,data)
  return data
}

func renameItem(args []string, data TodoListItems) TodoListItems {
  selectIndex,_ := strconv.Atoi(args[1])
  newTitle := args[2]
  data.Items[selectIndex].Title = newTitle
  listItems(args,data)
  return data
}

func toggleItem(args []string, data TodoListItems) TodoListItems {
  toggleIndex,_ := strconv.Atoi(args[1])
  data.Items[toggleIndex].Done = !data.Items[toggleIndex].Done
  listItems(args,data)
  return data
}
