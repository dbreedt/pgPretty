package interfaces

type SqlPrinter interface {
	AutoIndentAndPrintString(val string)
	PrintString(val string)
	PrintStringNoIndent(val string)
	PrintInt(val int)
	PrintInt64(val int64)
	PrintFloat64(val float64)
	PrintIntNoIndent(val int)
	PrintInt64NoIndent(val int64)
	PrintFloat64NoIndent(val float64)
	IncIndent()
	DecIndent()
	PrintKeywordNoIndent(keyword string)
	PrintKeyword(keyword string)
	NewLine()
	String() string
}
