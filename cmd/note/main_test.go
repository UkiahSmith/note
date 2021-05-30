package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func Test_run(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error("error getting current working directory: ", err)
	}

	testTmplPath := path.Join(wd, "testdata/19851026__review-slides.note")
	testGoldenPath := path.Join(wd, "testdata/19851026__review-slides__golden.md")

	var tmpDir = t.TempDir()
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Error(err)
	}

	// We are using `true` as our editor as a way to do non-interactive edit.
	// The template will be read, and the file created, so we can test, but no
	// editor will actually be run.
	var args = []string{"--date", "1985-10-26", "-e", "true", "-t", testTmplPath, "Review Slides"}

	// The golden version of the file we want.
	want, err := ioutil.ReadFile(testGoldenPath)
	if err != nil {
		t.Error(err)
		return
	}

	err = run(args, os.Stdout)
	if err != nil {
		fmt.Println(err)
		t.Error("error with run: ", err)
	}

	have, err := ioutil.ReadFile(path.Join(tmpDir, "19851026__review-slides.md"))
	if err != nil {
		t.Errorf("error reading created file: %s", err)
	}

	if bytes.Compare(want, have) != 0 {
		fmt.Println("bytes.Compare is ", bytes.Compare(want, have))
		fmt.Println("")
		fmt.Println("want: \n", string(want))
		fmt.Println("have: \n", string(have))
		fmt.Println("")
		t.Error("error: generated file does not match golden.")
	}
}
