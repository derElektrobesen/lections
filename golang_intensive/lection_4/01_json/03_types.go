package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type simpleType string

const (
	String simpleType = "string"
	Number simpleType = "number"
	Array  simpleType = "array"
	Object simpleType = "object"
)

func (t *simpleType) UnmarshalJSON(data []byte) error {
	supportedTypes := map[simpleType]bool{
		String: true,
		Number: true,
		Array:  true,
		Object: true,
	}

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid type field: %q: string required", data)
	}

	data = data[1 : len(data)-1]

	if possibleType := simpleType(data); supportedTypes[possibleType] {
		*t = possibleType
		return nil
	}

	return fmt.Errorf("unsupported data type: %q", data)
}

type item struct {
	Title       string     `json:"title"`
	Type        simpleType `json:"type"`
	Description string     `json:"description"`

	// integer type
	Minimum          *int  `json:"minimum"`
	ExclusiveMinumum *bool `json:"exclusiveMinimum"`

	// array type
	Items     *item `json:"items"`
	MinItems  *int  `json:"minItems"`
	UniqItems *bool `json:"uniqueItems"`

	// object type
	Properties map[string]item `json:"properties"`

	Required []string `json:"required"`
}

type schema struct {
	Schema string `json:"$schema"`

	item
}

func main() {
	var data schema
	if err := readJson(&data); err != nil {
		log.Fatalf("can't read json data: %s", err)
	}

	dataToPrint, _ := json.MarshalIndent(data, "", "  ")
	fmt.Printf("%s", dataToPrint)
}
