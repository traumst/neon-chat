package template

import (
	"bytes"
	"fmt"
	"html/template"
)

type ExpandableTemplate struct {
	Items []string
}

func (et ExpandableTemplate) HTML() (string, error) {
	if err := et.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("static/html/menu/expandable.html"))
	if err := tmpl.Execute(&buf, et); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (et ExpandableTemplate) validate() error {
	return nil
}
