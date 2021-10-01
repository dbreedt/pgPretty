package test

import (
	"reflect"
	"strings"
)

type Keywords struct {
	Ws        string
	Select    string
	With      string
	As        string
	From      string
	Limit     string
	On        string
	Where     string
	Join      string
	Group     string
	By        string
	Order     string
	Into      string
	Distinct  string
	Lateral   string
	Left      string
	Right     string
	Outer     string
	Count     string
	Over      string
	Partition string
	Lower     string
	Upper     string
	Maximum   string
	Minimum   string
	Any       string
	All       string
	Cross     string
	Full      string
}

func NewKeywords(upper bool) Keywords {
	retVal := Keywords{}

	val := reflect.Indirect(reflect.ValueOf(&retVal))
	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name

		if upper {
			name = strings.ToUpper(name)
		} else {
			name = strings.ToLower(name)
		}

		val.Field(i).SetString(name)
	}

	return retVal
}
