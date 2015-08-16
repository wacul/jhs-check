package schema

type Schema struct {
	Schema      string `json:"$schema"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Refer       string `json:"$ref"`
	Definitions map[string]Schema
	Properties  map[string]Schema
}

type Link struct {
	Title       string `json:"title"`
	Href        string `json:"href"`
	Method      string `json:"method"`
	Rel         string `json:"rel"`
	Description string `json:"description"`
	Schema      Schema `json:"schema"`
}
