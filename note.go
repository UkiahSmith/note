package note

import (
	"fmt"
	"io"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/gosimple/slug"
	"github.com/ukiahsmith/duolog"
)

var log = duolog.New(os.Stdout, "Note", 0)

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	ErrNotFound           = Err("Error: Not Found")
	ErrNoFilenameTemplate = Err("Error: No filename template found in first line")
)

type Data struct {
	Date    time.Time
	Title   string
	Content string
	Meta    struct {
		FilenameNOX string
	}
	tmpl *template.Template
}

func (d Data) TitleSlug() string {
	return slug.Make(d.Title)
}

// ParseDate validates a user supplied datetime against the three supported
// date formats, and set the date. It defaults to the current datetime if no
// datetime is supplied.
//
// 2006-01-02T15:04:05Z07:00
// 2006-01-02 15:04:05
// 2006-01-02
func (d *Data) ParseDate(usrDate string) error {
	var err error

	// Default to right about now, the funk soul brother.
	d.Date = time.Now()

	if usrDate != "" {
		var t time.Time
		t, err = time.Parse("2006-01-02T15:04:05Z07:00", usrDate)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04:05", usrDate)
			if err != nil {
				t, err = time.Parse("2006-01-02", usrDate)
				if err != nil {
					log.Infof("error parsing time format 2006-01-02: %s", err)
					return fmt.Errorf("error parsing date: %w", err)
				}
			}
		}

		// No errors? Use the user supplied date.
		if !t.IsZero() {
			d.Date = t
		}
	}

	return nil
}

func (d *Data) SetTemplateFile(templateFilename string) error {
	var err error

	// use the supplied template, or fallback to the default template
	if templateFilename == "" {
		d.tmpl, err = template.New("note").Funcs(Tfuncs).Parse(DefaultTmpl)
		if err != nil {
			return err
		}
	} else {
		d.Meta.FilenameNOX, err = FilenameFromFile(templateFilename, *d)
		log.Debugf("SetTemplateFile: using custom template with -t, d.Meta.FilenameNOX is \"%s\" after being set with note.FilenameFromFile", d.Meta.FilenameNOX)
		if err != nil {
			return err
		}

		d.tmpl, err = template.New(path.Base(templateFilename)).Funcs(Tfuncs).ParseFiles(templateFilename)
		if err != nil {
			return fmt.Errorf("%s: %w", err, ErrNotFound)
		}
	}

	if d.Meta.FilenameNOX == "" {
		d.Meta.FilenameNOX, err = FilenameFromTemplateStr(DefaultFilenameTmpl, *d)
		if err != nil {
			log.Debugf("error with default filenameTemplate: %s", err)
			return err
		}
	}

	return nil
}

func (d Data) GetFilename(filenameTemplate string) (string, error) {
	var fname string
	var err error

	// use the supplied filename-template, then fallback to a filename-template
	// in the note template, then fallback to the defaultly formatted filename.
	// var fname string

	// default names are setup when parsing the template.
	if filenameTemplate != "" {
		fname, err = FilenameFromTemplateStr(filenameTemplate, d)
		if err != nil {
			log.Debugf("error with default filenameTemplate: %s", err)
			return "", err
		}
		return fname, nil
	}

	return d.Meta.FilenameNOX, nil

}

func (d Data) Execute(wr io.Writer) error {
	return d.tmpl.Execute(wr, d)
}
