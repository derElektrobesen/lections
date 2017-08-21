package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func writeXML(data interface{}) error {
	fName := "test.xml"
	file, err := os.Create(fName)
	if err != nil {
		return fmt.Errorf("can't open XML file %q: %s", fName, err)
	}

	dec := xml.NewEncoder(file)
	dec.Indent("", "  ")
	if err := dec.Encode(data); err != nil {
		return fmt.Errorf("can't write XML data into file %q: %s", fName, err)
	}

	return nil
}

func readYAML(dest interface{}) error {
	fName := "test.yaml"

	file, err := os.Open(fName)
	if err != nil {
		return fmt.Errorf("can't open YAML file %q: %s", fName, err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("can't read yaml file %q: %s", fName, err)
	}

	if err := yaml.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("can't write XML data into file %q: %s", fName, err)
	}

	return nil
}

func main() {
	var cfg map[string]jigurdaConfig
	if err := readYAML(&cfg); err != nil {
		log.Fatalf("can't read jigurda config: %s", err)
	}

	var jigurda []jigurdaConfigXML
	for k, v := range cfg {
		var phones []string
		for k := range v.Phones {
			phones = append(phones, k)
		}

		j := jigurdaConfigXML{
			XMLName:     xml.Name{"", "jigurda"},
			JigurdaName: k,
			Permille:    v.Permille,
			Phones:      phones,
		}

		if v.MobileConfig != nil {
			j.MobileConfig = v.MobileConfig.toXML()
		}

		if c := v.WebConfig; c != nil {
			disabled := c.disableableConfig.toXML()
			j.WebConfig = &disabled
		}

		jigurda = append(jigurda, j)
	}

	err := writeXML(struct {
		XMLName  xml.Name
		Jigurdas []jigurdaConfigXML
	}{
		XMLName:  xml.Name{"", "jigurdas"},
		Jigurdas: jigurda,
	})

	if err != nil {
		log.Fatalf("can't write jigurda in XML file: %s", err)
	}
}
