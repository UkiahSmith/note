package note

import (
	"fmt"
	"os"
	"os/exec"
)

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
