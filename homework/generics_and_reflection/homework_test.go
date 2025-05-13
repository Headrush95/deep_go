package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(value any) string {
	obj := reflect.ValueOf(value)
	tp := reflect.TypeOf(value)

	sb := new(strings.Builder)
	countOfFields := obj.NumField()
	sb.Grow(countOfFields * 15) // чтобы хотя б не нулевой был

	for i := range countOfFields {
		if !tp.Field(i).IsExported() {
			continue
		}
		tag := tp.Field(i).Tag
		tagInfo, present := tag.Lookup("properties")
		if !present || len(tagInfo) == 0 || tagInfo == "-" {
			continue
		}
		splittedTag := strings.Split(tagInfo, ",")
		isOmitempty := false
		if len(splittedTag) == 2 && splittedTag[1] == "omitempty" {
			isOmitempty = true
		}
		switch obj.Field(i).Kind() {
		case reflect.String:
			if isOmitempty && obj.Field(i).IsZero() {
				continue
			}

			sb.WriteString(splittedTag[0] + "=" + obj.Field(i).String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if isOmitempty && obj.Field(i).IsZero() {
				continue
			}

			sb.WriteString(splittedTag[0] + "=" + strconv.FormatInt(obj.Field(i).Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if isOmitempty && obj.Field(i).IsZero() {
				continue
			}

			sb.WriteString(splittedTag[0] + "=" + strconv.FormatUint(obj.Field(i).Uint(), 10))
		case reflect.Float32, reflect.Float64:
			if isOmitempty && obj.Field(i).IsZero() {
				continue
			}

			sb.WriteString(splittedTag[0] + "=" + strconv.FormatFloat(obj.Field(i).Float(), 'g', -1, 64))
		case reflect.Bool:
			if isOmitempty && obj.Field(i).IsZero() {
				continue
			}

			sb.WriteString(splittedTag[0] + "=" + strconv.FormatBool(obj.Field(i).Bool()))
		case reflect.Array, reflect.Slice, reflect.Map:
			if isOmitempty && obj.Field(i).IsZero() {
				continue
			}

			sb.WriteString(splittedTag[0] + "=" + fmt.Sprintf("%v", obj.Field(i).Interface()))
		default:
			continue
		}

		if i < countOfFields-1 {
			sb.WriteString("\n")
		}
	}

	return strings.Clone(sb.String())
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
