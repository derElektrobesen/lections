package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func readJson(dest interface{}) error {
	fName := "test.json"
	file, err := os.Open(fName)
	if err != nil {
		return fmt.Errorf("can't open json file %q: %s", fName, err)
	}

	dec := json.NewDecoder(file)
	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("can't raad json data from file %q: %s", fName, err)
	}

	return nil
}
