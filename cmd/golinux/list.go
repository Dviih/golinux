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

