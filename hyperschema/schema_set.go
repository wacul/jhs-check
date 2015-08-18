package hyperschema

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type SchemaSet struct {
	schemas map[string]*Schema
	refs    map[string]string
	once    sync.Once
}

func (s *SchemaSet) init() {
	s.once.Do(func() {
		s.schemas = map[string]*Schema{}
		s.refs = map[string]string{}
	})
}

func (s *SchemaSet) AddYAML(path string) error {
	s.init()
	var raw RawSchema
	content, e := ioutil.ReadFile(path)
	if e != nil {
		return ErrorNotSupportedFile(e)
	}

	if e := yaml.Unmarshal(content, &raw); e != nil {
		return ErrorNotSupportedFile(e)
	}

	if raw.ID == "" {
		return ErrorNoID()
	}
	if strings.HasSuffix(raw.ID, "/") {
		return ErrorExtraSlash(raw.ID)
	}

	if raw.Schema == "" {
		return ErrorProperty(raw.ID, "schema", ErrorPropertyTypeEmpty)
	} else if raw.Schema != "http://json-schema.org/draft-04/hyper-schema" {
		return ErrorPropertyIncorrect(raw.ID, "schema", raw.Schema)
	}

	idPath := strings.Split(raw.ID, "/")
	if _, e := s.register(idPath[len(idPath)-1], &raw); e != nil {
		return e
	}
	return nil
}

func (s *SchemaSet) AddJSON(path string) error {
	s.init()
	//UNDONE:
	return nil
}

func (s *SchemaSet) Collect(path string, info os.FileInfo, e error) error {
	if e != nil {
		return e
	}
	if info.IsDir() {
		return nil
	}

	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		if e := s.AddYAML(path); e != nil {
			return e
		}
	} else if strings.HasSuffix(path, ".json") || strings.HasSuffix(path, ".js") {
		if e := s.AddJSON(path); e != nil {
			return e
		}
	} else {
		return ErrorNotSupportedFile(nil)
	}

	return nil
}

func (s *SchemaSet) Validate() error {
	var err *Error
	err = nil
	for id, ref := range s.refs {
		if ref == "" {
			continue
		}
		if !strings.HasPrefix(ref, "#/") {
			addError(&err, ErrorPropertyIncorrect(id, "$ref", ref))
			continue
		}
		refto := strings.TrimPrefix(ref, "#/")
		if _, ok := s.schemas[refto]; !ok {
			addError(&err, ErrorInvalidRefs(id, ref))
		}
	}
	//NOTE: return errとそのまま返すと、なぜかnil判定が失敗するのでここでnil判定しておく
	if err == nil {
		return nil
	} else {
		return err
	}
}

func (s *SchemaSet) register(id string, raw *RawSchema) (schema *Schema, err *Error) {
	if raw == nil {
		return
	}

	if id == "" {
		// 識別子のない登録は行われないはず
		panic("bug: empty id has entered in 'SchemaSet::register' method.")
	}

	if _, ok := s.schemas[id]; ok {
		addError(&err, ErrorDuplicated(id))
	}

	definitions, e := s.registerMap(id, "definitions", raw.Definitions)
	addError(&err, e)
	properties, e := s.registerMap(id, "properties", raw.Properties)
	addError(&err, e)
	links, e := s.registerLinks(id, raw.Links)
	addError(&err, e)
	items, e := s.register(id+"/items", raw.Items)
	addError(&err, e)

	var types []string
	if raw.Type == nil {
		//noop
	} else if typesArray, ok := raw.Type.([]interface{}); ok {
		types = make([]string, len(typesArray))
		for i, it := range typesArray {
			if st, ok := it.(string); ok {
				types[i] = st
			} else {
				addError(&err, ErrorPropertyIncorrect(id, "type", raw.Type))
				break
			}
		}
	} else if typeText, ok := raw.Type.(string); ok {
		types = []string{typeText}
	} else {
		addError(&err, ErrorPropertyIncorrect(id, "type", raw.Type))
	}

	var minValue float64
	if raw.Min != nil {
		if raw.MinValue != nil {
			addError(&err, ErrorProperty(id, "min, minValue", ErrorPropertyTypeConflicted))
		}
		minValue = *raw.Min
	} else if raw.MinValue != nil {
		minValue = *raw.MinValue
	}
	var maxValue float64
	if raw.Max != nil {
		if raw.MaxValue != nil {
			addError(&err, ErrorProperty(id, "max, maxValue", ErrorPropertyTypeConflicted))
		}
		maxValue = *raw.Max
	} else if raw.MaxValue != nil {
		maxValue = *raw.MaxValue
	}

	//NOTE: Refの内容チェックは、Validateメソッドで明示的に行う。

	if raw.Ref == "" && (types == nil || len(types) == 0) {
		addError(&err, ErrorProperty(id, "$ref", ErrorPropertyTypeNil))
		addError(&err, ErrorProperty(id, "type", ErrorPropertyTypeNil))
	}
	if raw.Ref != "" && types != nil && len(types) > 0 {
		addError(&err, ErrorProperty(id, "$ref, type", ErrorPropertyTypeConflicted))
	}

	for _, t := range types {
		switch t {
		case "object":
			// noop
			// if properties == nil {
			// 	addError(&err, ErrorProperty(id, "properties", ErrorPropertyTypeNil))
			// }
		case "array":
			if items == nil {
				addError(&err, ErrorProperty(id, "items", ErrorPropertyTypeNil))
			}
		case "string", "bool", "boolean", "integer", "number":
			// noop
		default:
			addError(&err, ErrorPropertyIncorrect(id, "type", t))
		}
	}

	schema = &Schema{
		id:     id,
		Schema: raw.Schema,

		Definitions: definitions,
		Links:       links,

		Ref: raw.Ref,

		Types:      types,
		Properties: properties,
		Pattern:    raw.Pattern,
		Enum:       raw.Enum,
		Items:      items,

		Required:  raw.Required,
		MinValue:  minValue,
		MaxValue:  maxValue,
		MinLength: raw.MinLength,
		MaxLength: raw.MaxLength,

		Description: raw.Description,
		Example:     raw.Example,
	}
	s.schemas[id] = schema
	s.refs[id] = raw.Ref
	return
}

func (s *SchemaSet) registerLinks(id string, rawLinks []*RawLink) (links []*Link, err *Error) {
	if rawLinks == nil {
		return
	}
	links = make([]*Link, len(rawLinks))
	for i := range rawLinks {
		child, e := s.register(fmt.Sprintf("%s/links/%d", id, i), rawLinks[i].Schema)
		addError(&err, e)
		links[i] = &Link{
			Title:       rawLinks[i].Title,
			Href:        rawLinks[i].Href,
			Method:      rawLinks[i].Method,
			Rel:         rawLinks[i].Rel,
			Description: rawLinks[i].Description,
			Schema:      child,
		}
	}

	return
}

func (s *SchemaSet) registerMap(id string, ref string, source map[string]*RawSchema) (schemas map[string]*Schema, err *Error) {
	if source == nil {
		return
	}

	var childID func(string) string
	if id == "" {
		childID = func(string) string {
			return ""
		}
	} else {
		childID = func(name string) string {
			return strings.Join([]string{id, ref, name}, "/")
		}
	}

	schemas = map[string]*Schema{}
	for name, def := range source {
		child, e := s.register(childID(name), def)
		addError(&err, e)
		schemas[name] = child
	}
	return
}

func (s *SchemaSet) Walk(walkFn func(s *Schema, err error) error) error {
	var err error
	for _, schema := range s.schemas {
		err = walkFn(schema, err)
		if err == SkipSchema {
			err = nil
		}
	}
	return err
}
