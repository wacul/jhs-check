package hyperschema

type RawSchema struct {
	ID string `json:"ID"`

	Schema      string                `json:"$schema" yaml:"$schema"`
	Type        interface{}           `json:"type"` // string or []string
	Description string                `json:"description"`
	Ref         string                `json:"$ref" yaml:"$ref"`
	Example     interface{}           `json:"example"`
	Pattern     string                `json:"pattern"`
	Required    []string              `json:"required"`
	Min         *float64              `json:"min"`
	MinValue    *float64              `json:"minValue"`
	Max         *float64              `json:"max"`
	MaxValue    *float64              `json:"maxValue"`
	MinLength   int                   `json:"minLength"`
	MaxLength   int                   `json:"maxLength"`
	Definitions map[string]*RawSchema `json:"definitions"`
	Properties  map[string]*RawSchema `json:"properties"`
	Enum        []string              `json:"enum"`
	Items       *RawSchema            `json:"items"`
	Links       []*RawLink            `json:"links"`
}

type RawLink struct {
	Title       string     `json:"title"`
	Href        string     `json:"href"`
	Method      string     `json:"method"`
	Rel         string     `json:"rel"`
	Description string     `json:"description"`
	Schema      *RawSchema `json:"schema"`
}
