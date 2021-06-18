package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

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
	if err := run(os.Args[1:], os.Stdout); err != nil {
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
	var filenameTemplate *string = fset.StringP("filename-format", "", "", "A valid Go template used to format the note's filename. e.g. {{ .TitleSlug }}.md")
	var editor *string = fset.StringP("editor", "e", "", "The text editor you want to open the generated note.")
	fset.StringVar(&noteD.Title, "title", "", "Use this to pre-populate the title variable in a template.")
	fset.StringVar(&noteD.Content, "content", "", "Use this to pre-populate the content variable in a template.")
	var tempDate *string = fset.String("date", "", "Use this to pre-populate the date variable in a template.")
	var showVersion *bool = fset.BoolP("version", "v", false, "Display the vesion information.")

	err = fset.Parse(args)
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

	err = noteD.ParseDate(*tempDate)
	if err != nil {
		fmt.Printf("error: date is invalid. %s is not parsable.\n", *tempDate)
		os.Exit(0)
	}

	var ed string
	if *editor != "" {
		ed = strings.TrimSpace(*editor)
	} else {
		ed = os.Getenv("EDITOR")
		if ed == "" {
			fset.Usage()
			fmt.Printf("\n\nerror: '-e' '--editor' or $EDITOR not set.\n")
			os.Exit(0)
		}
	}

	err = noteD.SetTemplateFile(*templateFile)
	if err != nil {
		fset.Usage()
		fmt.Printf("\n\nerror: %s\n", err)
		os.Exit(0)
	}

	fname, err := noteD.GetFilename(*filenameTemplate)
	if err != nil {
		fmt.Printf("\n\nerror: %s\n", err)
		os.Exit(0)
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
		runEditor(ed, fname)
		return nil
	}

	err = noteD.Execute(writer)
	if err != nil {
		fset.Usage()
		return fmt.Errorf("error executing template: %w", err)
	}

	writer.Close()

	err = runEditor(ed, fname)
	if err != nil {
		fset.Usage()
		return err
	}

	return nil
}

func runEditor(editor, filename string) error {
	path, err := exec.LookPath(editor)
	if err != nil {
		return fmt.Errorf("RunEditor: failed to find editor: %s", err)
	}

	cmd := exec.Command(path, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error running editor: %s", err)
	}

	return nil
}
