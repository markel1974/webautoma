package executor

import (
	"bytes"
	"github.com/google/uuid"
	"strings"
	"text/template"
)

// Templates represents a structure for managing and applying template variables to strings with custom mappings.
type Templates struct {
	variables map[string]interface{}
}

// NewTemplates initializes and returns a new Templates instance with the provided variables map.
func NewTemplates(variables map[string]interface{}) *Templates {
	return &Templates{
		variables: variables,
	}
}

// Apply processes a given string using predefined template variables and an additional stack of variables for replacements.
// Returns the processed string or an error if template parsing or execution fails.
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

// joinMaps merges two maps into a single map, giving precedence to the second map for duplicate keys.
// If both input maps are nil, it returns nil.
// If either map is nil, it returns the non-nil map.
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

// BuildCommand applies templates to the fields of a ConfigCommand using the given stack and updates the command fields.
func (t *Templates) BuildCommand(cmd ConfigCommand, stack map[string]interface{}) error {
	var err error
	if cmd.Id, err = t.Apply(cmd.Id, stack); err != nil {
		return err
	}
	if cmd.Command, err = t.Apply(cmd.Command, stack); err != nil {
		return err
	}
	if cmd.Target, err = t.Apply(cmd.Target, stack); err != nil {
		return err
	}
	if cmd.Until, err = t.Apply(cmd.Until, stack); err != nil {
		return err
	}
	if cmd.Value, err = t.Apply(cmd.Value, stack); err != nil {
		return err
	}
	return nil
}
