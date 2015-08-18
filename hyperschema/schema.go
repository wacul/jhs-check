package hyperschema

import "errors"

type Schema struct {
	id string

	Schema      string
	Types       []string
	Description string
	Ref         string
	Example     interface{}
	Pattern     string
	Required    []string
	MinValue    float64
	MaxValue    float64
	MinLength   int
	MaxLength   int
	Definitions map[string]*Schema
	Properties  map[string]*Schema
	Enum        []string
	Items       *Schema
	Links       []*Link
}

type Link struct {
	Title       string
	Href        string
	Method      string
	Rel         string
	Description string
	Schema      *Schema
}

var SkipSchema = errors.New("skip")

func (s *Schema) Walk(walkFn func(s *Schema, err error) error) error {
	if s == nil {
		return nil
	}
	err := walkFn(s, nil)
	if err == SkipSchema {
		return nil
	}
	for _, def := range s.Definitions {
		err = walkFn(def, err)
	}
	for _, prop := range s.Properties {
		err = walkFn(prop, err)
	}
	err = walkFn(s.Items, err)
	for _, link := range s.Links {
		err = walkFn(link.Schema, err)
	}
	return err
}
