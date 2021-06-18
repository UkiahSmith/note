package note

import (
	"testing"
	"time"
)

func Test_FilenameFromFile(t *testing.T) {
	date, err := time.Parse("20060102", "19851026")
	if err != nil {
		t.Fatalf("error initializing date: %s", err)
	}
	data := Data{Date: date, Title: "My test note "}

	out, err := FilenameFromFile("testdata/test_template.note", data)
	if err != nil {
		t.Fatalf("error with file: %s", err)
	}

	if err != nil || out != "19851026__my-test-note.md" {
		t.Errorf("FilenameFromFile failed: %w", err)
	}
}

func Test_FilenameFromFile_EmptyFilename(t *testing.T) {
	date, err := time.Parse("20060102", "19851026")
	if err != nil {
		t.Fatalf("error initializing date: %s", err)
	}
	data := Data{Date: date, Title: "My test note "}

	want := "my-test-note.md"
	got, err := FilenameFromFile("testdata/test_template-empty-filename.note", data)
	if err != nil {
		t.Fatalf("error with file: %s", err)
	}

	if err != nil || got != want {
		t.Errorf("FilenameFromFile failed, expected %s got %s : %s", want, got, err)
	}
}

func Test_FilenameFromFile_FileNotFound(t *testing.T) {
	date, err := time.Parse("20060102", "19851026")
	if err != nil {
		t.Fatalf("error initializing date: %s", err)
	}
	data := Data{Date: date, Title: "My test note "}

	_, err = FilenameFromFile("testdata/XXX__this-file-does-not-exist_52cf217a11ca67b2.note", data)
	if err != ErrNotFound {
		t.Fatalf("FilenameFromFile failed: expected file to be \"not found\", but error was: %+v", err)
	}
}

func Test_FilenameFromFile_CommentWithFilename(t *testing.T) {
	date, err := time.Parse("20060102", "19851026")
	if err != nil {
		t.Fatalf("error initializing date: %s", err)
	}
	data := Data{Date: date, Title: "My test note "}

	want := "19851026__my-test-note.md"
	got, err := FilenameFromFile("testdata/test_template-comment-with-filename.note", data)
	if err != nil {
		t.Fatalf("error with file: %s", err)
	}

	if err != nil || got != want {
		t.Errorf("FilenameFromFile failed, expected %s got %s : %s", want, got, err)
	}
}

func Test_FileFromFile_CommentWithFilename(t *testing.T) {
	date, err := time.Parse("20060102", "19851026")
	if err != nil {
		t.Fatalf("error initializing date: %s", err)
	}
	data := Data{Date: date, Title: "My test note "}

	want := "19851026__my-test-note.md"
	got, err := FilenameFromFile("testdata/test_template-comment-with-filename.note", data)
	if err != nil {
		t.Fatalf("error with file: %s", err)
	}

	if err != nil || got != want {
		t.Errorf("FilenameFromFile failed, expected %s got %s : %s", want, got, err)
	}
}

func Test_GetFirstLineFromTemplateFile(t *testing.T) {
	out, err := GetFirstLineFromTemplateFile("testdata/test_template.note")
	if err != nil {
		t.Errorf("GetFirstLineFromTemplateFile error: %s", err)
	}
	if out != `+++ #  {{ dateFormat "20060102" .Date }}__{{ makeSlug .Title }}.md` {
		t.Errorf("GetFirstLineFromTemplateFile failed.")
	}
}

func Test_GetFirstLineFromTemplateFile_FileNotFound(t *testing.T) {
	_, err := GetFirstLineFromTemplateFile("xxxtestdata/test_template_empty-filename.note")
	if err != ErrNotFound {
		t.Errorf("GetFirstLineFromTemplateFile expected file to be \"not found\", but error was: %+v", err)
	}
}

func Test_ExtractTemplateFromLine(t *testing.T) {
	out, err := ExtractTemplateFromLine(`+++ #  {{ dateFormat "20060102" .Date }}__{{ makeSlug .Title }}.md`)
	if err != nil {
		t.Errorf("ExtractTemplateFromLine error: %s", err)
	}

	if out != `{{ dateFormat "20060102" .Date }}__{{ makeSlug .Title }}.md` {
		t.Errorf("ExtractTemplateFromLine failed: %w", err)
	}
}

func Test_FilenameFromTemplate(t *testing.T) {
	date, err := time.Parse("20060102", "19851026")
	if err != nil {
		t.Fatalf("error initializing date: %s", err)
	}
	data := Data{Date: date, Title: "My test note "}
	out, err := FilenameFromTemplateStr(`{{ dateFormat "20060102" .Date }}__{{ makeSlug .Title }}.md`, data)
	if err != nil || out != "19851026__my-test-note.md" {
		t.Errorf("FilenameFromTemplate failed: %w", err)
	}
}
