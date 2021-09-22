package helpers

import (
	"reflect"
	"regexp"
)

/*
ProcessNamedParameters Scans the given sql and extracts any named parameters and replaces them with
	bare parameters.
	ORMs like go-pg/pg have support for named sql parameter placeholders like `?ClientID`, where
	Postgres sql syntax only supports bare `?` parameter place holders

	Note: This approach is very naive and does not support sql with a mix of named and bare parameters
*/
func ProcessNamedParameters(sql string) (string, map[int]string) {
	compRegEx := regexp.MustCompile(`(\?\w+)`)

	retVal := make(map[int]string)

	for i, match := range compRegEx.FindAllString(sql, -1) {
		retVal[i] = match
	}

	return compRegEx.ReplaceAllString(sql, "?"), retVal
}

// NilCheck Generic nil check
func NilCheck(n interface{}) bool {
	if n == nil {
		return true
	}

	switch reflect.TypeOf(n).Kind() {
	case reflect.Ptr,
		reflect.Map,
		reflect.Array,
		reflect.Chan,
		reflect.Slice:

		return reflect.ValueOf(n).IsNil()
	}

	return false
}
