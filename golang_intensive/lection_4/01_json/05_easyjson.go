package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/derElektrobesen/lections/golang_intensive/lection_4/01_json/test"
)

func main() {
	var data test.Schema
	if err := readJson(&data); err != nil {
		log.Fatalf("can't read json data: %s", err)
	}

	dataToPrint, _ := json.MarshalIndent(data, "", "  ")
	fmt.Printf("%s", dataToPrint)
}
