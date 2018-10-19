package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/UkiahSmith/note"
)

type NoteData struct {
	Date    time.Time
	Title   string
	Content string
}

func main() {
	var noteD NoteData
	var err error
	noteD.Date = time.Now()
	noteD.Title = strings.Join(os.Args[1:], " ")
	fname := strings.ToLower(strings.Join(os.Args[1:], "-")) + ".md"
	fmt.Println(noteD)
	fmt.Println(fname)

	ed := os.Getenv("EDITOR")
	if ed == "" {
		fmt.Println("error: $EDITOR not set.")
		os.Exit(1)
	}

	var writer io.WriteCloser
	_, err = os.Stat(fname)
	switch err.(type) {
	case *os.PathError:
		writer, err = os.Create(fname)
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}
	case error:
		fmt.Println("error: ", err)
		os.Exit(1)
	default:
		note.RunEditor(ed, fname)
		return
	}

	defer writer.Close()

	t, err := template.New("basic").Parse(basicNote)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	err = t.Execute(writer, noteD)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	note.RunEditor(ed, fname)
}

var basicNote = `+++
created_at = "{{ .Date }}"
modified_at = {{ .Date }}""
title = "{{ .Title }}"
+++

{{ .Content }}
`
