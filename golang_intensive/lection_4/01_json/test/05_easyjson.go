package test

type Item struct {
	Title       string `json:",omitempty"`
	Type        string `json:",omitempty"`
	Description string `json:",omitempty"`

	// integer type
	Minimum          *int  `json:",omitempty"`
	ExclusiveMinumum *bool `json:",omitempty"`

	// array type
	Items     *Item `json:",omitempty"`
	MinItems  *int  `json:",omitempty"`
	UniqItems *bool `json:",omitempty"`

	// object type
	Properties map[string]Item `json:",omitempty"`

	Required []string `json:",omitempty"`
}

type Schema struct {
	Schema string `json:",omitempty"`

	Item
}
