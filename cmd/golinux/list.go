package main

import (
	"github.com/Dviih/golinux/config"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"reflect"
	"strings"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Item struct {
	title       interface{}
	description interface{}
}

func (item *Item) Title() string {
	switch title := item.title.(type) {
	case string:
		return title
	case reflect.Value:
		return toString(title)
	case func() string:
		return title()
	default:
		return "<invalid title>"
	}
}

func (item *Item) Description() string {
	switch description := item.description.(type) {
	case string:
		return description
	case reflect.Value:
		return toString(description)
	case func() string:
		return description()
	default:
		return "<invalid description>"
	}
}

func (item *Item) FilterValue() string {
	return item.Title()
}

type backItem struct{}

func (item *backItem) Title() string {
	return "Back"
}

func (item *backItem) Description() string {
	return "Go back"
}

func (item *backItem) FilterValue() string {
	return ""
}

type createItem struct {
	t reflect.Type
}

func (item *createItem) Title() string {
	return "Create"
}

func (item *createItem) Description() string {
	return "Create " + item.t.String()
}

func (item *createItem) FilterValue() string {
	return ""
}

type List struct {
	Title        string
	size         tea.WindowSizeMsg
	list         list.Model
	rv           reflect.Value
	currentInput *ListInput
	back         *List
}


func (l *List) Init() tea.Cmd {
	rv := rvAbs(l.rv)
	delegate := list.NewDefaultDelegate()

	var items []list.Item

	if rv.Type() == reflect.TypeFor[config.KVS]() {
		for i := 0; i < rv.Len(); i++ {
			kv := rv.Index(i).Elem()

			items = append(items, &Item{
				title:       kv.Field(0),
				description: kv.Field(1),
			})
		}
	} else {
		switch rv.Kind() {
		case reflect.Map:
			m := rv.MapRange()

			for m.Next() {
				var description interface{}

				mv := rvAbs(m.Value())

				switch mv.Kind() {
				case reflect.Map:
					description = "Enter to access Map"
				case reflect.Slice:
					description = "Enter to access slice"
				case reflect.Struct:
					description = "Enter to access struct"
				default:
					key := m.Key()
					description = func() string {
						return toString(rv.MapIndex(key))
					}
				}

				items = append(items, &Item{
					title:       m.Key(),
					description: description,
				})
			}

			items = append(items, &createItem{t: rv.Type().Elem()})
		case reflect.Slice:
			delegate.ShowDescription = false

			for i := 0; i < rv.Len(); i++ {
				element := rvAbs(rv.Index(i))

				var title interface{}

				switch element.Kind() {
				case reflect.Map:
					title = "Enter to access Map"
				case reflect.Slice:
					title = "Enter to access slice"
				case reflect.Struct:
					title = "Enter to access struct"
				default:
					title = element
				}

				items = append(items, &Item{
					title: title,
				})
			}
		case reflect.Struct:
			rt := rv.Type()

			for i := 0; i < rt.NumField(); i++ {
				sf := rt.Field(i)

				if !sf.IsExported() {
					continue
				}

				field := rv.Field(i)

				var description interface{}

				switch field.Kind() {
				case reflect.Map:
					description = "Enter to access map"
				case reflect.Slice:
					description = "Enter to access slice"
				case reflect.Struct:
					description = "Enter to access struct"
				default:
					description = rv.Field(i)
				}

				items = append(items, &Item{
					title:       sf.Name,
					description: description,
				})
			}
		}
	}

	if l.back != nil {
		items = append(items, &backItem{})
	}

	w, h := docStyle.GetFrameSize()

	w = l.size.Width - w
	h = l.size.Height - h

	if w < 0 {
		w = 0
	}

	if h < 0 {
		h = 0
	}

	l.list = list.New(items, delegate, w, h-10)
	l.list.DisableQuitKeybindings()
	l.list.SetFilteringEnabled(false)
	l.list.Title = l.title()

	return nil
}

