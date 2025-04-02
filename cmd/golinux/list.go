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

func (l *List) title() string {
	s := ""

	tmp := l

	for tmp != nil {
		s = strings.ToUpper(string(tmp.Title[0])) + tmp.Title[1:] + "." + s
		tmp = tmp.back
	}

	return s[:len(s)-1]
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

func (l *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return l, tea.Quit
		case tea.KeyCtrlR:
			if l.currentInput != nil {
				break
			}

			rv := rvAbs(l.rv)

			if rv.Kind() != reflect.Map {
				break
			}

			item, ok := l.list.Items()[l.list.Index()].(*Item)
			if !ok {
				return l, nil
			}

			key := item.title.(reflect.Value)

			input := textinput.New()
			input.SetValue(toString(key))
			input.Focus()

			l.currentInput = &ListInput{
				action: ListInputActionRename,
				rv:     key,
				name:   toString(key),
				input:  input,
			}
		case tea.KeyEnter:
			if _, ok := l.list.Items()[l.list.Index()].(*backItem); ok {
				*l = *l.back
				return l, l.Init()
			}

			if l.currentInput != nil {
				if l.currentInput.input.Value() == "" && l.currentInput.action != ListInputActionCreate {
					l.currentInput = nil
					return l, nil
				}

				if l.currentInput.action == ListInputActionRename {
					rv := rvAbs(l.rv)

					nk := reflect.New(rv.Type().Key()).Elem()
					setValue(nk, l.currentInput.input.Value())

					rv.SetMapIndex(nk, rv.MapIndex(l.currentInput.rv))
					rv.SetMapIndex(l.currentInput.rv, reflect.Value{})

					l.currentInput = nil
					return l, l.Init()
				}

				switch m := l.currentInput.extra.(type) {
				case *LIMap:
					ptr := reflect.New(m.Get().Type()).Elem()

					setValue(ptr, l.currentInput.input.Value())
					m.Set(ptr)

					l.currentInput = nil
					return l, nil
				case *LICreate:
					switch l.rv.Kind() {
					case reflect.Map:
						key := reflect.New(l.rv.Type().Key()).Elem()
						setValue(key, l.currentInput.input.Value())

						ptr := reflect.New(m.t)

						for tmp := ptr.Elem(); tmp.Kind() == reflect.Pointer; {
							if !tmp.CanSet() {
								break
							}

							tmp.Set(reflect.New(tmp.Type().Elem()))
							tmp = tmp.Elem()
						}

						l.rv.SetMapIndex(key, ptr.Elem())
					}

					l.currentInput = nil
					return l, l.Init()
				default:
					setValue(l.currentInput.rv, l.currentInput.input.Value())
					l.currentInput = nil

					return l, nil
				}
			}

			if l.rv.Type() == reflect.TypeFor[config.KVS]() {
				kv := l.rv.Index(l.list.Index()).Elem()
				input := textinput.New()

				input.Focus()
				input.SetValue(toString(kv.Field(1)))

				l.currentInput = &ListInput{
					action: ListInputActionEdit,
					rv:     kv.Field(1),
					name:   toString(kv.Field(0)),
					input:  input,
				}

				goto end
			}

			if item, ok := l.list.Items()[l.list.Index()].(*createItem); ok {
				input := textinput.New()

				input.Focus()

				l.currentInput = &ListInput{
					action: ListInputActionCreate,
					rv:     l.rv,
					name:   item.t.String(),
					input:  input,
					extra:  &LICreate{t: item.t},
				}

				goto end
			}

			rv := rvAbs(l.rv)

			switch rv.Kind() {
			case reflect.Map:
				key := l.list.Items()[l.list.Index()].(*Item).title.(reflect.Value)
				mv := rvAbs(rv.MapIndex(key))

				switch mv.Kind() {
				case reflect.Map, reflect.Slice, reflect.Struct:
					back := *l

					m := &List{
						Title:        toString(key),
						size:         l.size,
						rv:           mv,
						currentInput: nil,
						back:         &back,
					}

					return m, m.Init()
				default:
					input := textinput.New()

					input.Focus()
					input.SetValue(toString(mv))

					l.currentInput = &ListInput{
						action: ListInputActionEdit,
						rv:     mv,
						name:   toString(key),
						input:  input,
						extra: &LIMap{
							m:   rv,
							key: key,
						},
					}
				}
			case reflect.Struct:
				rt := rv.Type()
				k := 0

				var (
					field reflect.Value
					ft    reflect.StructField
				)

				for i := 0; i < rt.NumField(); i++ {
					if !rt.Field(i).IsExported() {
						continue
					}

					if k == l.list.Index() {
						field = rv.Field(i)
						ft = rt.Field(i)
						break
					}

					k++
				}

				switch field.Kind() {
				case reflect.Map, reflect.Slice, reflect.Struct:
					back := *l

					m := &List{
						Title:        ft.Name,
						size:         l.size,
						rv:           field,
						currentInput: nil,
						back:         &back,
					}

					return m, m.Init()
				default:
					input := textinput.New()

					input.SetValue(toString(field))
					input.Focus()

					l.currentInput = &ListInput{
						action: ListInputActionEdit,
						rv:     field,
						name:   ft.Name,
						input:  input,
					}
				}
			}
		}
	case tea.WindowSizeMsg:
		w, h := docStyle.GetFrameSize()
		l.list.SetSize(msg.Width-w, msg.Height-h-10)
		l.size = msg
	}

end:
	if l.currentInput != nil {
		m, cmd := l.currentInput.input.Update(msg)
		l.currentInput.input = m

		return l, cmd
	}

	m, cmd := l.list.Update(msg)
	l.list = m

	return l, cmd
}

func (l *List) View() string {
	if l.currentInput != nil {
		switch extra := l.currentInput.extra.(type) {
		case *LIKV:
			return l.list.Styles.Title.Render(l.currentInput.Action()+l.title()) + "\n" + l.currentInput.input.View() + "\n" + extra.input.View()
		default:
			return l.list.Styles.Title.Render(l.currentInput.Action()+l.title()) + "\n" + l.currentInput.input.View()
		}
	}

	return docStyle.Render(l.list.View())
}

