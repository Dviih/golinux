package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"strings"
)

type Writer struct {
	program *tea.Program
}

func (writer *Writer) Write(data []byte) (int, error) {
	writer.program.Send(ExecData{Data: data})
	return len(data), nil
}

type ExecData struct {
	Data []byte
}

type ExecProgram struct {
	Program *tea.Program
}

