package note

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ukiahsmith/duolog"
)

var log = duolog.New(os.Stdout, "Note", 0)

type Data struct {
	Date      time.Time
	Title     string
	TitleSlug string
	Content   string
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
