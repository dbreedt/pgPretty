package test

import (
	"reflect"
	"strings"
)

type Keywords struct {
	FnUpper   bool
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
	And       string
	Not       string
	Between   string
	Or        string
	Like      string
	Is        string
	Null      string // not to sure about this move, null is value and not a keyword but people that write SELECT, FROM, etc expect NULL
	In        string
	Exists    string
	Desc      string
	Nulls     string
	Last      string
	Case      string
	Filter    string
	Having    string
	Ilike     string
}

func NewKeywords(upperCaseKeywords, upperCaseFunctions bool) Keywords {
	retVal := Keywords{}

	val := reflect.Indirect(reflect.ValueOf(&retVal))
	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name

		if name == "FnUpper" {
			val.Field(i).SetBool(upperCaseFunctions)
			continue
		}

		if upperCaseKeywords {
			name = strings.ToUpper(name)
		} else {
			name = strings.ToLower(name)
		}

		val.Field(i).SetString(name)
	}

	return retVal
}

func (k Keywords) Fn(name string) string {
	if k.FnUpper {
		return strings.ToUpper(name)
	}

	return strings.ToLower(name)
}
