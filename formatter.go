package main

import (
	"fmt"
	"strings"
)

type Formatter struct {
	indentation   string
	currentIndent int
	keywordInCaps bool
}

// NewFormatter Creates a new Formatter
func NewFormatter(tabs, keywordInCaps bool, numIndentations int) *Formatter {
	retVal := &Formatter{
		currentIndent: 0,
		keywordInCaps: keywordInCaps,
	}

	if tabs {
		retVal.indentation = strings.Repeat("\t", numIndentations)
	} else {
		retVal.indentation = strings.Repeat(" ", numIndentations)
	}

	return retVal
}

func (f *Formatter) makeIndent() string {
	return strings.Repeat(f.indentation, f.currentIndent)
}

// TODO: use a string buffer

// AutoIndentAndPrintString Manages indentation and prints values
func (f *Formatter) AutoIndentAndPrintString(val string) {
	f.currentIndent++
	fmt.Printf("%s%s", f.makeIndent(), val)
	f.currentIndent--
}

// PrintString Prints a string on the current indentation level
func (f *Formatter) PrintString(val string) {
	fmt.Printf("%s%s", f.makeIndent(), val)
}

// PrintStringNoIndent Prints a string
func (f *Formatter) PrintStringNoIndent(val string) {
	fmt.Printf("%s", val)
}

// PrintInt Prints an int on the current indentation level
func (f *Formatter) PrintInt(val int64) {
	fmt.Printf("%s%d", f.makeIndent(), val)
}

// PrintIntNoIndent Prints an int on the current indentation level
func (f *Formatter) PrintIntNoIndent(val int64) {
	fmt.Printf("%d", val)
}

// IncIndent Manually control the indent
func (f *Formatter) IncIndent() {
	f.currentIndent++
}

// DecIndent Manually control the indent
func (f *Formatter) DecIndent() {
	f.currentIndent--

	if f.currentIndent < 0 {
		panic("indent < 0")
	}
}

// PrintKeywordNoIndent Print a keyword with no indent
func (f *Formatter) PrintKeywordNoIndent(keyword string) {
	if f.keywordInCaps {
		fmt.Printf("%s", strings.ToUpper(keyword))
	} else {
		fmt.Printf("%s", keyword)
	}
}

// PrintKeyword Print a keyword on the current indentation level
func (f *Formatter) PrintKeyword(keyword string) {
	indent := f.makeIndent()

	if f.keywordInCaps {
		fmt.Printf("%s%s", indent, strings.ToUpper(keyword))
	} else {
		fmt.Printf("%s%s", indent, keyword)
	}
}

// NewLine prints a newline
func (f *Formatter) NewLine() {
	fmt.Printf("\n")
}
