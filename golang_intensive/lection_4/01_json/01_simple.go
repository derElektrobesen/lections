package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type item struct {
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`

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
	Title  string `json:"title"`
	Type   string `json:"type"`
	Items  item   `json:"items"`
}

func main() {
	var data schema
	if err := readJson(&data); err != nil {
		log.Fatalf("can't read json data: %s", err)
	}

	dataToPrint, _ := json.MarshalIndent(data, "", "  ")
	fmt.Printf("%s", dataToPrint)
}
