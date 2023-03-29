package executor

import (
	"bytes"
	"github.com/google/uuid"
	"strings"
	"text/template"
)

type Templates struct {
	variables map[string]interface{}
}

func NewTemplates(variables map[string]interface{}) *Templates {
	return &Templates{
		variables: variables,
	}
}

func (t *Templates) Apply(in string) (string, error) {
	if t.variables == nil {
		return in, nil
	}
	if !strings.Contains(in, "{{") {
		return in, nil
	}
	if !strings.Contains(in, "}}") {
		return in, nil
	}
	tmpl, err := template.New(uuid.New().String()).Parse(in)
	if err != nil {
		return "", err
	}
	out := bytes.NewBuffer(nil)
	err = tmpl.Execute(out, t.variables)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
