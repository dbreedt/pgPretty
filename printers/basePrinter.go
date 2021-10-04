package printers

import (
	"fmt"
	"strings"
)

type caseFormatter func(s string) string

type BasePrinter struct {
	indentation       string
	currentIndent     int
	keywordFormatter  caseFormatter
	functionFormatter caseFormatter
	sb                strings.Builder
	indentCache       map[int]string
}

// NewBasePrinter Creates a custom BasePrinter
func NewBasePrinter(tabs, keywordInCaps, functionInCaps bool, numIndentations int) *BasePrinter {
	retVal := &BasePrinter{
		currentIndent:     0,
		keywordFormatter:  strings.ToLower,
		functionFormatter: strings.ToLower,
		indentCache:       make(map[int]string, 5),
	}

	if keywordInCaps {
		retVal.keywordFormatter = strings.ToUpper
	}

	if functionInCaps {
		retVal.functionFormatter = strings.ToUpper
	}

	if tabs {
		retVal.indentation = strings.Repeat("\t", numIndentations)
	} else {
		retVal.indentation = strings.Repeat(" ", numIndentations)
	}

	return retVal
}

// NewSpacePrinter Create a semi generic BasePrinter that uses spaces
func NewSpacePrinter(keywordInCaps, functionInCaps bool, numIndentations int) *BasePrinter {
	return NewBasePrinter(false, keywordInCaps, functionInCaps, numIndentations)
}

// NewDefaultSpacePrinter Create a semi generic BasePrinter that uses spaces, lowercase keywords and lowercase function names
func NewDefaultSpacePrinter() *BasePrinter {
	return NewBasePrinter(false, false, false, 2)
}

// NewTabPrinter Create a semi generic BasePrinter that uses tabs
func NewTabPrinter(keywordInCaps, functionInCaps bool, numIndentations int) *BasePrinter {
	return NewBasePrinter(true, keywordInCaps, functionInCaps, numIndentations)
}

// NewDefaultTabPrinter Create a semi generic BasePrinter that uses tabs, lowercase keywords and lowercase function names
func NewDefaultTabPrinter() *BasePrinter {
	return NewBasePrinter(true, false, false, 1)
}

func (bp *BasePrinter) makeIndent() string {
	if v, ok := bp.indentCache[bp.currentIndent]; ok {
		return v
	}

	v := strings.Repeat(bp.indentation, bp.currentIndent)
	bp.indentCache[bp.currentIndent] = v
	return v
}

func (bp *BasePrinter) PrintString(val string, withIndent ...bool) {
	if len(withIndent) > 0 && withIndent[0] {
		bp.sb.WriteString(bp.makeIndent())
	}
	bp.sb.WriteString(val)
}

func (bp *BasePrinter) PrintInt(val int, withIndent ...bool) {
	bp.PrintInt64(int64(val), withIndent...)
}

func (bp *BasePrinter) PrintInt64(val int64, withIndent ...bool) {
	if len(withIndent) > 0 && withIndent[0] {
		bp.sb.WriteString(bp.makeIndent())
	}
	bp.sb.WriteString(fmt.Sprintf("%d", val))
}

func (bp *BasePrinter) PrintFloat64(val float64, withIndent ...bool) {
	if len(withIndent) > 0 && withIndent[0] {
		bp.sb.WriteString(bp.makeIndent())
	}
	bp.sb.WriteString(fmt.Sprintf("%f", val))
}

func (bp *BasePrinter) IncIndent() {
	bp.currentIndent++
}

func (bp *BasePrinter) DecIndent() {
	bp.currentIndent--

	if bp.currentIndent < 0 {
		panic("indent < 0")
	}
}

func (bp *BasePrinter) PrintKeyword(keyword string, withIndent ...bool) {
	if len(withIndent) > 0 && withIndent[0] {
		bp.sb.WriteString(bp.makeIndent())
	}
	bp.sb.WriteString(bp.keywordFormatter(keyword))
}

func (bp *BasePrinter) PrintFunction(functionName string, withIndent ...bool) {
	if len(withIndent) > 0 && withIndent[0] {
		bp.sb.WriteString(bp.makeIndent())
	}
	bp.sb.WriteString(bp.functionFormatter(functionName))
}

func (bp *BasePrinter) NewLine() {
	bp.sb.WriteString("\n")
}

func (bp *BasePrinter) String() string {
	return bp.sb.String()
}
