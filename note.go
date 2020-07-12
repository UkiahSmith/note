package note

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/gosimple/slug"
	"github.com/ukiahsmith/duolog"
)

var log = duolog.New(os.Stdout, "Note", 0)

type Data struct {
	Date    time.Time
	Title   string
	Content string
}

func (d Data) TitleSlug() string {
	return slug.Make(d.Title)
}

func RunEditor(editor, filename string) {
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
