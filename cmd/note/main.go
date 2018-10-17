package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

type Note struct {
	Date    time.Time
	Title   string
	Content string
}

func main() {
	var note Note
	note.Date = time.Now()
	note.Title = strings.Join(os.Args[1:], " ")
	fname := strings.ToLower(strings.Join(os.Args[1:], "-")) + ".md"
	fmt.Println(note)
	fmt.Println(fname)

	ed := os.Getenv("EDITOR")
	if ed == "" {
		fmt.Println("$EDITOR not set.")
		os.Exit(1)
	}

	var writer io.WriteCloser
	writer, err := os.Create(fname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer writer.Close()

	t, err := template.New("basic").Parse(basicNote)
	if err != nil {
		panic(err)
	}

	err = t.Execute(writer, note)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(ed, fname)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

}

var basicNote = `+++
created_at = "{{ .Date }}"
modified_at = {{ .Date }}""
title = "{{ .Title }}"
+++

{{ .Content }}
`
