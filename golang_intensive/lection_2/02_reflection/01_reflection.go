package main

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
)

//`some_tag:"set_to(100500)"`

type user struct {
	Name    string
	Surname string
	Age     int
}

var users = []user{
	user{
		Name:    "Vasya",
		Surname: "Pupkin",
		Age:     10,
	},
	user{
		Name:    "Petya",
		Surname: "Yakovlev",
		Age:     23,
	},
	user{
		Name:    "Ivan",
		Surname: "Petrov",
		Age:     17,
	},
	user{
		Name:    "Kirill",
		Surname: "Ivanov",
		Age:     32,
	},
}

func func1() {
	for i, u := range users {
		fmt.Printf("User %d: %s %s, age %d\n", i, u.Name, u.Surname, u.Age)
	}
}

func iterateStructure(data reflect.Value, cb func(string, reflect.Value)) {
	for i := 0; i < data.NumField(); i++ {
		fieldType := data.Type().Field(i)
		field := data.Field(i)

		cb(fieldType.Name, field)

		tagValue := fieldType.Tag.Get("some_tag")
		if tagValue == "" {
			continue
		}

		re := regexp.MustCompile(`^(.+)\((.+)\)$`)
		args := re.FindStringSubmatch(tagValue)
		if len(args) != 3 {
			log.Fatalf("invalid tag %q, args are %+v", tagValue, args)
		}

		funcName, funcArgs := args[1], args[2]

		if funcName != "set_to" {
			log.Fatalf("unsupported function %q", funcName)
		}

		if !field.CanSet() {
			log.Fatalf("Can't set field %q", fieldType.Name)
		}

		switch field.Kind() {
		case reflect.Int:
			i, _ := strconv.ParseInt(funcArgs, 10, 64)
			field.SetInt(i)
		case reflect.String:
			field.SetString(funcArgs)
		default:
			log.Fatalf("Unsupported kind to set field %q: %q", fieldType.Name, field.Kind())
		}

		cb(fieldType.Name+" (after modification)", field)
	}
}

func iterateArray(data reflect.Value, cb func(int, reflect.Value)) {
	for i := 0; i < data.Len(); i++ {
		cb(i, data.Index(i))
	}
}

func printTypesImpl(data reflect.Value, indent string) {
	prnt := func(format string, args ...interface{}) {
		args = append([]interface{}{indent}, args...)
		fmt.Printf("%s"+format+"\n", args...)
	}

	switch data.Kind() {
	case reflect.Interface:
		printTypesImpl(reflect.ValueOf(data.Interface()), indent)
	case reflect.Int:
		prnt("Integer found: %d", data.Int())
	case reflect.String:
		prnt("String found: %q", data.String())
	case reflect.Struct:
		prnt("Structure of type %q found", data.Type())
		iterateStructure(data, func(fieldName string, data reflect.Value) {
			printTypesImpl(data, fmt.Sprintf("%s  %s: ", indent, fieldName))
		})
	case reflect.Array:
	case reflect.Slice:
		prnt("Array of type %q with len %d found", data.Type(), data.Len())
		iterateArray(data, func(i int, data reflect.Value) {
			printTypesImpl(data, fmt.Sprintf("%s  %d: ", indent, i))
		})
	case reflect.Func:
		prnt("Function of type %q found", data.Type())
		retArgs := data.Call([]reflect.Value{})
		for _, arg := range retArgs {
			printTypesImpl(arg, fmt.Sprintf("%s  ret: ", indent))
		}
	default:
		prnt("Unexpected type found: %q", data.Type())
	}
}

func printTypes(data interface{}) {
	printTypesImpl(
		reflect.ValueOf(data),
		"",
	)
}

func func2() {
	printTypes(users)
	printTypes([]interface{}{
		1,
		"xxxx",
		struct {
			X int
			Y string
			Z int64
			K struct {
				L int
			}
		}{},
		[]interface{}{
			[]byte{},
			"yyyyy",
		},
		func() string { return "some string" },
	})
}

func main() {
	func1()
	//func2()
}
