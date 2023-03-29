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

func (t *Templates) Apply(in string, stack map[string]interface{}) (string, error) {
	if !strings.Contains(in, "{{") {
		return in, nil
	}
	if !strings.Contains(in, "}}") {
		return in, nil
	}

	computed := t.joinMaps(t.variables, stack)
	if computed == nil {
		return in, nil
	}
	tmpl, err := template.New(uuid.New().String()).Parse(in)
	if err != nil {
		return "", err
	}
	out := bytes.NewBuffer(nil)
	err = tmpl.Execute(out, computed)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func (t *Templates) joinMaps(in1 map[string]interface{}, in2 map[string]interface{}) map[string]interface{} {
	if in1 == nil && in2 == nil {
		return nil
	}
	if in1 != nil && in2 == nil {
		return in1
	}
	if in2 != nil && in1 == nil {
		return in2
	}
	joined := make(map[string]interface{})
	if in1 != nil {
		for k, v := range in1 {
			joined[k] = v
		}
	}
	if in2 != nil {
		for k, v := range in2 {
			joined[k] = v
		}
	}
	return joined
}
