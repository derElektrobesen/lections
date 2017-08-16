package main

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"strings"
)

type enumType []string

type property struct {
	Type        string   `json:"type"`
	Enum        enumType `json:"enum"`
	Description string   `json:"description"`
	Name        string   `json:"-"`
}

type schema struct {
	Type           string               `json:"type"`
	Properties     map[string]*property `json:"properties"`
	RequiredFields []string             `json:"required"`
}

func main() {
	var sch schema
	if err := json.Unmarshal(templateBaseData, &sch); err != nil {
		log.Fatalf("Can't unmarshal schema: %s", err)
	}

	for k, v := range sch.Properties {
		if v.Type == "" {
			// for enums
			v.Type = strings.Title(k) + "Type"
		}
	}

	tmpl.Execute(os.Stdout, sch)
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")

	var ret string
	for _, p := range parts {
		ret += strings.Title(p)
	}

	return ret
}

var templateBaseData = []byte(`
{
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "description": "Имя пользователя"
    },
    "surname": {
      "type": "string",
      "description": "Фамилия пользователя"
    },
    "birth_date": {
      "type": "string",
      "description": "Дата рождения"
    },
    "gender": {
      "enum": [
        "male",
        "female"
      ],
      "description": "Пол пользователя"
    },
    "about": {
      "type": "string",
      "description": "Информация о пользователе"
    },
    "start_driving_date": {
      "type": "string",
      "description": "дата, когда юзер сел за руль"
    },
    "email": {
      "type": "string",
      "description": "Email адрес пользователя"
    }
  },
  "required": [
    "name",
    "surname",
    "birth_date",
    "gender",
    "email"
  ],
  "$schema": "http://json-schema.org/draft-04/schema#"
}`)

var tmpl = template.Must(template.New("").
	Funcs(template.FuncMap{
		"title": strings.Title,
		"camel": toCamelCase,
		"ap":    func() string { return "`" },
	}).
	Parse(`
package profileschema

{{- range $key, $value := .Properties }}
	{{- if ne (len $value.Enum) 0 }}
		type {{ $value.Type }} string // {{ $value.Description }}
		const (
			{{- range $value.Enum }}
			{{ title . }} {{ title $key }}Type = "{{ . }}"
			{{- end }}
		)
	{{- end }}
{{- end }}

type Profile struct {
{{- range $key, $value := .Properties }}
	{{- $t := $value.Type }}
	{{ camel $key }} {{ $t }} {{ ap }}json:"{{ $key }}"{{ ap }} // {{ $value.Description }}
{{- end }}
}

`))
