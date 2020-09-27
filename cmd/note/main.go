package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/ukiahsmith/duolog"
	"github.com/ukiahsmith/note"
)

const (
	exitFail = 1
)

var (
	buildVersion   string = "dev"
	buildTimestamp string = ""
	buildHash      string = ""
)

var log = duolog.New(os.Stderr, "Note", 0)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		log.Info(err)
		os.Exit(exitFail)
	}
}

func run(args []string, stdout io.Writer) error {
	var noteD note.Data
	var fset flag.FlagSet
	var err error

	fset.Usage = func() {
		fmt.Print(`
Note, a templating tool for note taking.

Usage:
	note [options] <Title of note>

Options:

`)
		fset.PrintDefaults()
		fmt.Print(`
Note:
	The title flag _or_ title arguments is required.

	The date flag must be in one of these formats
		2006-01-02T15:04:05Z07:00
		2006-01-02 15:04:05
		2006-01-02

`)
	}

	var templateFile *string = fset.StringP("template", "t", "", "The file to use as a template")
	var filenameTemplate *string = fset.StringP("filename", "", "", "A valid Go template used to format the note's filename. e.g. {{ .TitleSlug }}.md")
	fset.StringVar(&noteD.Title, "title", "", "Use this to pre-populate the title variable in a template.")
	fset.StringVar(&noteD.Content, "content", "", "Use this to pre-populate the content variable in a template.")
	var tempDate *string = fset.String("date", "", "Use this to pre-populate the date variable in a template.")
	var showVersion *bool = fset.BoolP("version", "v", false, "Display the vesion information.")

	err = fset.Parse(args[1:])
	if err != nil {
		log.Infof("error parsing arguments: %v", err)
		os.Exit(1)
	}

	if *showVersion {
		fmt.Println("Note, a templating tool for note taking.")
		fmt.Println("Version:    ", buildVersion)
		if !strings.HasPrefix(buildVersion, "dev") {
			fmt.Println("Build date: ", buildTimestamp)
			fmt.Println("Build hash: ", buildHash)
		}
		os.Exit(0)
	}

	if noteD.Title == "" {
		noteD.Title = strings.TrimSpace(strings.Join(fset.Args(), " "))
	}

	if noteD.Title == "" {
		fset.Usage()
		return errors.New("error: Title is required")
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
					if err != nil {
						log.Infof("error parsing time format 2006-01-02: %v", err)
					}
				}
			}
			if !t.IsZero() {
				noteD.Date = t
			}
		}
	}

	ed := os.Getenv("EDITOR")
	if ed == "" {
		fset.Usage()
		return errors.New("error: $EDITOR not set.")
	}

	// use the supplied template, or fallback to the default template
	var tmpl *template.Template
	if *templateFile == "" {
		tmpl, err = template.New("note").Funcs(note.Tfuncs).Parse(note.DefaultTmpl)
	} else {
		tmpl, err = template.New(*templateFile).Funcs(note.Tfuncs).ParseFiles(*templateFile)
	}
	if err != nil {
		fset.Usage()
		return fmt.Errorf("error finding template: %w", err)
	}

	// use the supplied filename-template, then fallback to a filename-template
	// in the note template, then fallback to the defaultly formatted filename.
	var fname string
	switch {
	case *filenameTemplate != "":
		fname, err = note.FilenameFromTemplateStr(*filenameTemplate, noteD)
		if err != nil {
			fset.Usage()
			log.Debugf("error with default filenameTemplate: %s", err)
			return err
		}
	case *templateFile != "":
		if *templateFile != "" {
			// returns a filename using the note.DefaultFilenameTmpl if there is no valid filename-template in the note-template.
			fname = note.FilenameFromFile(*templateFile, noteD)
			log.Debugf("got to using custom template with -t, fname is \"%s\" after being set with note.FilenameFromFile", fname)
		}
	default:
		fname, err = note.FilenameFromTemplateStr(note.DefaultFilenameTmpl, noteD)
		if err != nil {
			fset.Usage()
			log.Debugf("error with default filenameTemplate: %s", err)
			return err
		}
	}

	var writer io.WriteCloser
	_, err = os.Stat(fname)
	switch err.(type) {
	case *os.PathError:
		writer, err = os.Create(fname)
		if err != nil {
			fset.Usage()
			return fmt.Errorf("os.PathError with file %s : %w", fname, err)
		}
	case error:
		fset.Usage()
		return fmt.Errorf("error with the file: %w", err)
	default:
		note.RunEditor(ed, fname)
		return nil
	}

	defer writer.Close()

	err = tmpl.Execute(writer, noteD)
	if err != nil {
		fset.Usage()
		return fmt.Errorf("error executing template: %w", err)
	}

	note.RunEditor(ed, fname)
	fmt.Println("Created file: ", fname)

	return nil
}
