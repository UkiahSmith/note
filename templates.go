package note

import (
	"text/template"
	"time"
)

var BasicTmpl = `+++
title = "{{ .Title }}"
created_at = "{{ dateFormat "2006-01-02T15:04:05Z07:00" .Date }}"
modified_at = "{{ dateFormat "2006-01-02T15:04:05Z07:00" .Date }}"
+++

{{ .Content }}
`

func DateFormat(layout string, t time.Time) string {
	return t.Format(layout)
}

var Tfuncs template.FuncMap = map[string]interface{}{
	"dateFormat": DateFormat,
}
