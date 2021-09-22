package printers

import (
	"fmt"
	"strings"
)

type caseFormatter func(s string) string

type BasePrinter struct {
	indentation      string
	currentIndent    int
	keywordFormatter caseFormatter
	sb               strings.Builder
	indentCache      map[int]string
}

// NewBasePrinter Creates a custom BasePrinter
func NewBasePrinter(tabs, keywordInCaps bool, numIndentations int) *BasePrinter {
	retVal := &BasePrinter{
		currentIndent:    0,
		keywordFormatter: strings.ToLower,
		indentCache:      make(map[int]string, 5),
	}

	if keywordInCaps {
		retVal.keywordFormatter = strings.ToUpper
	}

	if tabs {
		retVal.indentation = strings.Repeat("\t", numIndentations)
	} else {
		retVal.indentation = strings.Repeat(" ", numIndentations)
	}

	return retVal
}

// NewSpacePrinter Create a semi generic BasePrinter that uses spaces
func NewSpacePrinter(keywordInCaps bool, numIndentations int) *BasePrinter {
	return NewBasePrinter(false, keywordInCaps, numIndentations)
}

// NewDefaultSpacePrinter Create a semi generic BasePrinter that uses spaces
func NewDefaultSpacePrinter() *BasePrinter {
	return NewBasePrinter(false, false, 2)
}

// NewTabPrinter Create a semi generic BasePrinter that uses tabs
func NewTabPrinter(keywordInCaps bool, numIndentations int) *BasePrinter {
	return NewBasePrinter(true, keywordInCaps, numIndentations)
}

// NewDefaultTabPrinter Create a semi generic BasePrinter that uses spaces
func NewDefaultTabPrinter() *BasePrinter {
	return NewBasePrinter(true, false, 2)
}

func (bp *BasePrinter) makeIndent() string {
	if v, ok := bp.indentCache[bp.currentIndent]; ok {
		return v
	}

	v := strings.Repeat(bp.indentation, bp.currentIndent)
	bp.indentCache[bp.currentIndent] = v
	return v
}

// AutoIndentAndPrintString Manages indentation and prints values
func (bp *BasePrinter) AutoIndentAndPrintString(val string) {
	bp.currentIndent++
	bp.sb.WriteString(bp.makeIndent())
	bp.sb.WriteString(val)
	bp.currentIndent--
}

// PrintString Prints a string on the current indentation level
func (bp *BasePrinter) PrintString(val string) {
	bp.sb.WriteString(bp.makeIndent())
	bp.PrintStringNoIndent(val)
}

// PrintStringNoIndent Prints a string
func (bp *BasePrinter) PrintStringNoIndent(val string) {
	bp.sb.WriteString(val)
}

// PrintInt Prints an int on the current indentation level
func (bp *BasePrinter) PrintInt(val int) {
	bp.PrintInt64(int64(val))
}

func (bp *BasePrinter) PrintInt64(val int64) {
	bp.sb.WriteString(bp.makeIndent())
	bp.PrintInt64NoIndent(val)
}

// PrintIntNoIndent Prints an int on the current indentation level
func (bp *BasePrinter) PrintIntNoIndent(val int) {
	bp.PrintInt64NoIndent(int64(val))
}

func (bp *BasePrinter) PrintInt64NoIndent(val int64) {
	bp.sb.WriteString(fmt.Sprintf("%d", val))
}

func (bp *BasePrinter) PrintFloat64(val float64) {
	bp.sb.WriteString(bp.makeIndent())
	bp.PrintFloat64NoIndent(val)
}

func (bp *BasePrinter) PrintFloat64NoIndent(val float64) {
	bp.sb.WriteString(fmt.Sprintf("%f", val))
}

// IncIndent Manually control the indent
func (bp *BasePrinter) IncIndent() {
	bp.currentIndent++
}

// DecIndent Manually control the indent
func (bp *BasePrinter) DecIndent() {
	bp.currentIndent--

	if bp.currentIndent < 0 {
		panic("indent < 0")
	}
}

// PrintKeywordNoIndent Print a keyword with no indent
func (bp *BasePrinter) PrintKeywordNoIndent(keyword string) {
	bp.sb.WriteString(bp.keywordFormatter(keyword))
}

// PrintKeyword Print a keyword on the current indentation level
func (bp *BasePrinter) PrintKeyword(keyword string) {
	bp.sb.WriteString(bp.makeIndent())
	bp.PrintKeywordNoIndent(keyword)
}

// NewLine prints a newline
func (bp *BasePrinter) NewLine() {
	bp.sb.WriteString("\n")
}

func (bp *BasePrinter) String() string {
	return bp.sb.String()
}
