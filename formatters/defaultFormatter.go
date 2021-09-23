package formatters

import (
	"fmt"
	"strconv"

	interfaces "github.com/dbreedt/pgPretty/interfaces"
	nodes "github.com/pganalyze/pg_query_go/nodes"
)

type DefaultFormatter struct {
	printer            interfaces.SqlPrinter
	detectedParameters map[int]string
	paramCounter       int
}

func NewDefaultFormatterWithParameters(printer interfaces.SqlPrinter, parameterLookup map[int]string) *DefaultFormatter {
	return &DefaultFormatter{
		printer:            printer,
		detectedParameters: parameterLookup,
	}
}

func NewDefaultFormatter(printer interfaces.SqlPrinter) *DefaultFormatter {
	return NewDefaultFormatterWithParameters(printer, nil)
}

func (df *DefaultFormatter) p(msg string) {
	// Dump the printer's content before we panic to aid with debugging
	fmt.Println(df.printer)

	panic(msg + " not supported")
}

func (df *DefaultFormatter) PrintWithClause(wc nodes.WithClause) {
	if wc.Recursive {
		df.p("With Clause - Recursive")
	}

	for i := range wc.Ctes.Items {
		if i == 0 {
			df.printer.PrintKeywordNoIndent("with ")
		}

		df.PrintNode(wc.Ctes.Items[i], false)

		if i < len(wc.Ctes.Items)-1 {
			df.printer.PrintString(", ")
		}
	}
}

func (df *DefaultFormatter) PrintSelectStatement(ss nodes.SelectStmt) {
	if ss.WithClause != nil {
		df.PrintWithClause(*ss.WithClause)
	}

	df.printer.PrintKeyword("select")
	df.printer.NewLine()
	df.printer.IncIndent()

	for i := range ss.DistinctClause.Items {
		if ss.DistinctClause.Items[i] == nil {
			df.printer.PrintKeywordNoIndent("distinct")
			df.printer.NewLine()
		} else {
			df.PrintNode(ss.DistinctClause.Items[i], true)
		}
	}

	if ss.IntoClause != nil {
		df.p("Select - Into clause")
	}

	for i := range ss.TargetList.Items {
		df.PrintNode(ss.TargetList.Items[i], true)

		if i < len(ss.TargetList.Items)-1 {
			df.printer.PrintStringNoIndent(",")
			df.printer.NewLine()
		}
	}

	df.printer.DecIndent()
	df.printer.NewLine()

	for i := range ss.FromClause.Items {
		if je, ok := ss.FromClause.Items[i].(nodes.JoinExpr); ok {
			df.PrintJoin(i == 0, je)
		} else {
			df.printer.PrintKeyword("from")
			df.printer.NewLine()
			df.printer.IncIndent()
			df.PrintNode(ss.FromClause.Items[i], true)
		}

		df.printer.NewLine()
	}

	df.printer.DecIndent()
	df.printer.PrintKeyword("where")
	df.printer.NewLine()
	df.printer.IncIndent()
	df.PrintNode(ss.WhereClause, true)
	df.printer.DecIndent()
}

func (df *DefaultFormatter) PrintResTarget(nt nodes.ResTarget, withIndent bool) {
	retVal := ""

	if nt.Name != nil {
		retVal = *nt.Name
	}

	for i := range nt.Indirection.Items {
		df.PrintNode(nt.Indirection.Items[i], withIndent)
	}

	df.PrintNode(nt.Val, withIndent)

	df.printer.PrintStringNoIndent(retVal)
}

func (df *DefaultFormatter) PrintColumnRef(cr nodes.ColumnRef, withIndent bool) {
	for i := range cr.Fields.Items {
		df.PrintNode(cr.Fields.Items[i], withIndent)

		if i < len(cr.Fields.Items)-1 {
			df.printer.PrintStringNoIndent(".")
		}

		// negate indentation on consecutive runs to stop tab.<space>col prints
		withIndent = false
	}
}

func (df *DefaultFormatter) PrintJoin(first bool, join nodes.JoinExpr) {
	if first {
		df.printer.PrintKeyword("from")
		df.printer.NewLine()
		df.printer.IncIndent()
	} else {
		df.printer.IncIndent()
		df.PrintJoinType(join.Jointype, true)
	}

	if join.IsNatural {
		df.p("Join - Natural")
	}

	// cross join
	if join.Jointype == nodes.JOIN_INNER && join.Quals == nil {
		df.PrintNode(join.Larg, true)
		df.printer.NewLine()
		df.printer.DecIndent()
		df.printer.PrintKeyword("cross join ")
		df.PrintNode(join.Rarg, false)
		df.printer.IncIndent()
	} else {
		df.PrintNode(join.Larg, true)
		df.printer.NewLine()
		df.printer.DecIndent()
		df.PrintJoinType(join.Jointype, true)
		df.PrintNode(join.Rarg, false)

		if len(join.UsingClause.Items) > 0 {
			df.p("Join - Using Clause")
		}

		df.printer.NewLine()
		df.printer.IncIndent()
		df.printer.PrintKeyword("on ")
		df.PrintNode(join.Quals, false)
	}
}

func (df *DefaultFormatter) PrintAlias(alias *nodes.Alias) {
	if alias == nil {
		return
	}

	if alias.Aliasname != nil {
		df.printer.PrintStringNoIndent(*(alias.Aliasname))
	}

	if len(alias.Colnames.Items) > 0 {
		df.p("Alias - Columns")
	}
}

func (df *DefaultFormatter) PrintJoinType(joinType nodes.JoinType, withIndent bool) {
	jt := ""

	switch joinType {
	case nodes.JOIN_INNER:
		jt = "join "

	case nodes.JOIN_LEFT:
		jt = "left join "

	case nodes.JOIN_FULL:
		jt = "full join "

	case nodes.JOIN_RIGHT:
		jt = "right join "

	case nodes.JOIN_SEMI:
		df.p("join type - semi")

	case nodes.JOIN_ANTI:
		df.p("join type - anti")

	case nodes.JOIN_UNIQUE_OUTER:
		df.p("join type - unique outer")

	case nodes.JOIN_UNIQUE_INNER:
		df.p("join type - unique inner")
	}

	if len(jt) > 0 {
		if withIndent {
			df.printer.PrintKeyword(jt)
		} else {
			df.printer.PrintKeywordNoIndent(jt)
		}
	}
}

func (df *DefaultFormatter) PrintRangeVar(rv nodes.RangeVar, withIndent bool) {
	name := ""

	if rv.Catalogname != nil {
		name = *rv.Catalogname
	}

	if rv.Schemaname != nil {
		if len(name) > 0 {
			name += "."
		}

		name += *rv.Schemaname
	}

	if rv.Relname != nil {
		if len(name) > 0 {
			name += "."
		}

		name += *rv.Relname
	}

	if withIndent {
		df.printer.PrintString(name)
	} else {
		df.printer.PrintStringNoIndent(name)
	}

	if rv.Alias != nil {
		df.printer.PrintStringNoIndent(" ")
		df.PrintAlias(rv.Alias)
	}
}

func (df *DefaultFormatter) PrintBoolExpr(be nodes.BoolExpr, prevOp nodes.BoolExprType, withIndent bool) {
	if be.Xpr != nil {
		df.PrintNode(be.Xpr, withIndent)
		df.printer.PrintStringNoIndent(" ")
		df.PrintBoolExprType(be.Boolop, false)
	}

	// Drop parentheses when change operators
	parentheses := be.Boolop != prevOp && be.Boolop != 2

	if parentheses {
		if withIndent {
			df.printer.PrintString("(")
		} else {
			df.printer.PrintStringNoIndent("(")
		}

		df.printer.NewLine()
		df.printer.IncIndent()
	}

	for i := range be.Args.Items {
		if int(be.Boolop) == 2 {
			df.PrintBoolExprType(be.Boolop, withIndent)
		}

		if tbe, ok := be.Args.Items[i].(nodes.BoolExpr); ok {
			df.PrintBoolExpr(tbe, be.Boolop, withIndent)
		} else {
			// only print an indent if this is the first item and the operator is not a not
			df.PrintNode(be.Args.Items[i], i == 0 && (int(be.Boolop) != 2))
		}

		if i < len(be.Args.Items)-1 {
			df.printer.NewLine()
			df.PrintBoolExprType(be.Boolop, true)
		}

		// Stop printing consecutive args with indent to stop crap like AND<space>tab.col...
		withIndent = false
	}

	if parentheses {
		df.printer.NewLine()
		df.printer.DecIndent()
		df.printer.PrintString(")")
	}
}

func (df *DefaultFormatter) PrintBoolExprType(exprType nodes.BoolExprType, withIndent bool) {
	kw := ""

	switch exprType {
	case nodes.AND_EXPR:
		kw = "and "

	case nodes.OR_EXPR:
		kw = "or "

	default:
		if int(exprType) == 2 {
			kw = "not "
		}
	}

	if len(kw) > 0 {
		if withIndent {
			df.printer.PrintKeyword(kw)
		} else {
			df.printer.PrintKeywordNoIndent(kw)
		}
	}
}

func (df *DefaultFormatter) PrintAExpr(ae nodes.A_Expr, withIndent bool) {
	df.PrintNode(ae.Lexpr, withIndent)

	switch ae.Kind {
	case nodes.AEXPR_OP_ANY:
		df.PrintAExprKeywords(ae.Name)
		df.printer.PrintKeywordNoIndent("any")
		df.printer.PrintStringNoIndent("( ")
		df.PrintNode(ae.Rexpr, false)
		df.printer.PrintStringNoIndent(" )")

	case nodes.AEXPR_OP_ALL:
		df.PrintAExprKeywords(ae.Name)
		df.printer.PrintKeywordNoIndent("all")
		df.printer.PrintStringNoIndent("( ")
		df.PrintNode(ae.Rexpr, false)
		df.printer.PrintStringNoIndent(" )")

	case nodes.AEXPR_BETWEEN,
		nodes.AEXPR_NOT_BETWEEN:

		df.PrintAExprKeywords(ae.Name)
		df.PrintNode(ae.Rexpr.(nodes.List).Items[0], false)
		df.printer.PrintKeywordNoIndent(" and ")
		df.PrintNode(ae.Rexpr.(nodes.List).Items[1], false)

	case nodes.AEXPR_LIKE:
		df.printer.PrintKeywordNoIndent(" like ")
		df.PrintNode(ae.Rexpr, false)

	case nodes.AEXPR_ILIKE:
		df.printer.PrintKeywordNoIndent(" ilike ")
		df.PrintNode(ae.Rexpr, false)

	case nodes.AEXPR_OP:
		// Prints the operator
		df.PrintAExprKeywords(ae.Name)

		df.PrintNode(ae.Rexpr, false)

	default:
		df.p(fmt.Sprintf("AExpr: %v", ae.Kind))
	}
}

func (df *DefaultFormatter) PrintAExprKeywords(op nodes.List) {
	for i := range op.Items {
		// These are keywords
		df.printer.PrintStringNoIndent(" ")
		df.printer.PrintKeywordNoIndent(op.Items[i].(nodes.String).Str)
		df.printer.PrintStringNoIndent(" ")
	}
}

func (df *DefaultFormatter) PrintNullTest(nt nodes.NullTest, withIndent bool) {
	if nt.Xpr != nil {
		df.PrintNode(nt.Xpr, withIndent)
	}

	df.PrintNode(nt.Arg, withIndent)

	if nt.Nulltesttype == nodes.IS_NULL {
		df.printer.PrintKeywordNoIndent(" is null")
	} else {
		df.printer.PrintKeywordNoIndent(" is not null")
	}
}

func (df *DefaultFormatter) PrintAConst(ac nodes.A_Const) {
	switch ac.Val.(type) {
	case nodes.String:
		df.printer.PrintStringNoIndent("'")
		df.PrintNode(ac.Val, false)
		df.printer.PrintStringNoIndent("'")

	default:
		df.PrintNode(ac.Val, false)
	}
}

func (df *DefaultFormatter) PrintCommonTableExpr(cte nodes.CommonTableExpr) {
	if cte.Ctename != nil {
		df.printer.PrintStringNoIndent(*cte.Ctename)
		df.printer.PrintStringNoIndent(" (")
		df.printer.NewLine()
		df.printer.IncIndent()
	}

	if cte.Cterecursive {
		df.p("CTE - Recursive")
	}

	df.PrintNode(cte.Ctequery, false)

	df.printer.NewLine()
	df.printer.DecIndent()
	df.printer.PrintString(")")
	df.printer.NewLine()
}

func (df *DefaultFormatter) PrintParamRef(pr nodes.ParamRef, withIndent bool) {
	if df.detectedParameters != nil {
		if param, ok := df.detectedParameters[df.paramCounter]; ok {
			if withIndent {
				df.printer.PrintString(param)
			} else {
				df.printer.PrintStringNoIndent(param)
			}
			df.paramCounter++
			return
		}
	}

	if withIndent {
		df.printer.PrintString("?")
	} else {
		df.printer.PrintStringNoIndent("?")
	}
}

func (df *DefaultFormatter) PrintNode(node nodes.Node, withIndent bool) {
	switch node.(type) {
	case nodes.RawStmt:
		df.PrintNode(node.(nodes.RawStmt).Stmt, withIndent)

	case nodes.SelectStmt:
		df.PrintSelectStatement(node.(nodes.SelectStmt))

	case nodes.IntoClause:
		df.printer.PrintKeywordNoIndent(" into ")

	case nodes.ResTarget:
		df.PrintResTarget(node.(nodes.ResTarget), withIndent)

	case nodes.ColumnRef:
		df.PrintColumnRef(node.(nodes.ColumnRef), withIndent)

	case nodes.A_Star:
		df.printer.PrintString("*")

	case nodes.String:
		if withIndent {
			df.printer.PrintString(node.(nodes.String).Str)
		} else {
			df.printer.PrintStringNoIndent(node.(nodes.String).Str)
		}

	case nodes.JoinExpr:
		df.PrintJoin(false, node.(nodes.JoinExpr))

	case nodes.RangeVar:
		df.PrintRangeVar(node.(nodes.RangeVar), withIndent)

	case nodes.BoolExpr:
		df.PrintBoolExpr(node.(nodes.BoolExpr), node.(nodes.BoolExpr).Boolop, withIndent)

	case nodes.A_Expr:
		df.PrintAExpr(node.(nodes.A_Expr), withIndent)

	case nodes.A_Const:
		df.PrintAConst(node.(nodes.A_Const))

	case nodes.Integer:
		if withIndent {
			df.printer.PrintInt64(node.(nodes.Integer).Ival)
		} else {
			df.printer.PrintInt64NoIndent(node.(nodes.Integer).Ival)
		}

	case nodes.Float:
		val, err := strconv.ParseFloat(node.(nodes.Float).Str, 64)
		if err != nil {
			df.p(node.(nodes.Float).Str + " failed to parse as a float")
		}

		if withIndent {
			df.printer.PrintFloat64(val)
		} else {
			df.printer.PrintFloat64NoIndent(val)
		}

	case nodes.NullTest:
		df.PrintNullTest(node.(nodes.NullTest), withIndent)

	case nodes.CommonTableExpr:
		df.PrintCommonTableExpr(node.(nodes.CommonTableExpr))

	case nodes.ParamRef:
		df.PrintParamRef(node.(nodes.ParamRef), withIndent)

	case nodes.RangeSubselect:
		// join lateral (select * from x where y = t.id)
		df.p("Node: Range Subselect")

	case nodes.SubLink:
		// and t.id in (select id from x where y)
		df.p("Node: Sublink")

	default:
		df.p(fmt.Sprintf("Node: %T", node))
	}
}