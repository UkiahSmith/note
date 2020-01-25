[![Build Status](https://travis-ci.com/UkiahSmith/note.svg?branch=master)](https://travis-ci.com/UkiahSmith/note) [![Go Report Card](https://goreportcard.com/badge/github.com/ukiahsmith/note)](https://goreportcard.com/report/github.com/ukiahsmith/note)

# Note

Note is a templating tool to enable quick creation of notes with useful meta information.

## Documentation

Note makes it easier to take notes by reading a template file, populating it with runtime data, and opening the result in an editor. It is useful to create templates for specific types of notes, such as internal meeting notes, client meetings, task notes, or even Readme.md that would quickly be populated with useful information with ease.

1. Install note as you would any Golang cli application.
2. Ensure that your `$EDITOR` environmental variable is set.
3. Run note using `note This is my first note`
4. note will take the arguments and use them as the Title of the note, and in templates.
5. The default template will be used. 

### Usage

```
Usage:
        note [options] <Title of note>

Options:

  -t, --template string   The file to use as a template
      --title string      Use this to pre-populate the title variable in a template.
      --slug string       Use this to pre-populate the slug variable.
      --content string    Use this to pre-populate the content variable in a template.
      --date string       Use this to pre-populate the date variable in a template.

Note:
        The title flag _or_ title arguments is required.

        The date flag must be in one of these formats
                2006-01-02T15:04:05Z07:00
                2006-01-02 15:04:05
                2006-01-02
```


### Default template

The default template is built into note itself, and will be used unless another template is specified on the command line by using the `-t` or `--template` flag.

```
+++
title = "{{ .Title }}"
created_at = "{{ dateFormat "2006-01-02T15:04:05Z07:00" .Date }}"
modified_at = "{{ dateFormat "2006-01-02T15:04:05Z07:00" .Date }}"
+++

{{ .Content }}
```

### Creating your own template

It is straight forward enough to create your own template. Note uses the [Golang templates](https://golang.org/pkg/text/template/) format. 

1. Create a text file, and use any of the available template values
2. Run `note -t mytemplate.txt` 

**Example custom template**

```
Title: {{ .Title }}
Date: {{ dateFormat "2 January 2006" .Date }}

{{ .Content }}
```

Most of the time the `.Content` value will be empty, however it is still useful to include it for the rare times when text is passed in with the `--content` flag.

#### Available template values

`.Title` is the title of the note, either from the arguments, or input from the `--title` flag.
`.TitleSlug` is the title normalized to a URL safe slug. 
`.Date` is a Go time.Time data type.
`.Content` is any content passed to note by using the `--content` flag.


#### Available template functions

There are functions available in a template to help transform data.

`dateFormat` is a function that takes a format string, a time.Time data type, and outputs the date as the format string.

```
{{ dateFormat "2006-01-02T15:04:05Z07:00" .Date }}

{{ dateFormat "2006-01-02" .Date }}

{{ dateFormat "Mon Jan 2 15:04:05 -0700 MST 2006" .Date }}
```

See the [Golang time.Format documentation](https://golang.org/pkg/time/#Time.Format) for details on the specifics of the format string.
