package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/UkiahSmith/note"
	"github.com/gosimple/slug"
	flag "github.com/spf13/pflag"
)

type NoteData struct {
	Date      time.Time
	Title     string
	TitleSlug string
	Content   string
}

func main() {
	var noteD NoteData
	var fset flag.FlagSet
	var err error

	fset.Usage = func() {
		fmt.Println(`
Note is a templating tool for note taking.

Usage:
	note [options] <Title of note>

Options:
`)
		fset.PrintDefaults()
		fmt.Println(`
Note:
	The title flag _or_ title arguments is required.

	The date flag must be in one of these formats
		2006-01-02T15:04:05Z07:00
		2006-01-02 15:04:05
		2006-01-02
`)
	}

	var templateFile *string = fset.StringP("template", "t", "", "The file to use as a template")
	fset.StringVar(&noteD.Title, "title", "", "Use this to pre-populate the title variable in a template.")
	fset.StringVar(&noteD.TitleSlug, "slug", "", "Use this to pre-populate the slug variable.")
	fset.StringVar(&noteD.Content, "content", "", "Use this to pre-populate the content variable in a template.")
	var tempDate *string = fset.String("date", "", "Use this to pre-populate the date variable in a template.")

	fset.Parse(os.Args[1:])

	if noteD.Title == "" {
		noteD.Title = strings.TrimSpace(strings.Join(fset.Args(), " "))
	}

	if noteD.Title == "" {
		fset.Usage()
		fmt.Println("\nerror: Title is required")
		os.Exit(1)
	}

	if noteD.TitleSlug == "" {
		noteD.TitleSlug = slug.Make(noteD.Title)
	} else {
		noteD.TitleSlug = slug.Make(noteD.TitleSlug)
	}

	{
		noteD.Date = time.Now()
		if *tempDate != "" {
			var t time.Time
			t, err = time.Parse("2006-01-02T15:04:05Z07:00", *tempDate)
			if err != nil {
				t, err = time.Parse("2006-01-02 15:04:05", *tempDate)
				if err != nil {
					t, err = time.Parse("2006-01-02", *tempDate)
				}
			}
			if !t.IsZero() {
				noteD.Date = t
			}
		}
	}

	fname := noteD.TitleSlug + ".md"

	ed := os.Getenv("EDITOR")
	if ed == "" {
		fset.Usage()
		fmt.Println("error: $EDITOR not set.")
		os.Exit(1)
	}

	var writer io.WriteCloser
	_, err = os.Stat(fname)
	switch err.(type) {
	case *os.PathError:
		writer, err = os.Create(fname)
		if err != nil {
			fset.Usage()
			fmt.Println("error: ", err)
			os.Exit(1)
		}
	case error:
		fset.Usage()
		fmt.Println("error: ", err)
		os.Exit(1)
	default:
		note.RunEditor(ed, fname)
		return
	}

	defer writer.Close()

	var tmpl *template.Template
	if *templateFile == "" {
		tmpl, err = template.New("note").Funcs(note.Tfuncs).Parse(note.BasicTmpl)
	} else {
		tmpl, err = template.New(*templateFile).Funcs(note.Tfuncs).ParseFiles(*templateFile)
	}
	if err != nil {
		fset.Usage()
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	err = tmpl.Execute(writer, noteD)
	if err != nil {
		fset.Usage()
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	note.RunEditor(ed, fname)
	fmt.Println("Created file: ", fname)
}
