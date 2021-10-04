package interfaces

type SqlPrinter interface {
	PrintString(val string, withIndent ...bool)
	PrintInt(val int, withIndent ...bool)
	PrintInt64(val int64, withIndent ...bool)
	PrintFloat64(val float64, withIndent ...bool)
	PrintKeyword(keyword string, withIndent ...bool)
	IncIndent()
	DecIndent()
	NewLine()
	String() string
}
