package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type simpleType interface {
	Type() string
}

type simpleTypeWrapper struct {
	Data simpleType
}

type item struct {
	Title       string `json:"title"`
	Description string `json:"description"`

	simpleTypeWrapper

	Required []string `json:"required"`
}

func (w *simpleTypeWrapper) UnmarshalJSON(data []byte) error {
	var t struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return fmt.Errorf("can't recognise item type: %s", err)
	}

	if t.Type == "" {
		// to `type` field found
		return nil
	}

	for _, supportedType := range []simpleType{
		&stringType{},
		&integerType{},
		&arrayType{},
		&objectType{},
	} {
		if supportedType.Type() == t.Type {
			w.Data = supportedType
			break
		}
	}

	if w.Data == nil {
		return fmt.Errorf("unsupportd type %q found", t.Type)
	}

	err := json.Unmarshal(data, w.Data)
	return err
}

type stringType struct{}

type integerType struct {
	Minimum          int  `json:"minimum"`
	ExclusiveMinumum bool `json:"exclusiveMinimum"`
}

type arrayType struct {
	Items     item `json:"items"`
	MinItems  int  `json:"minItems"`
	UniqItems bool `json:"uniqueItems"`
}

type objectType struct {
	Properties map[string]item `json:"properties"`
}

func (stringType) Type() string  { return "string" }
func (integerType) Type() string { return "number" }
func (arrayType) Type() string   { return "array" }
func (objectType) Type() string  { return "object" }

type schema struct {
	Schema string `json:"$schema"`

	item
}

func main() {
	var data schema
	if err := readJson(&data); err != nil {
		log.Fatalf("can't read json data: %s", err)
	}

	log.Printf(">>>>>>>>>> %+v", data)
	dataToPrint, _ := json.MarshalIndent(data, "", "  ")
	fmt.Printf("%s", dataToPrint)
}
