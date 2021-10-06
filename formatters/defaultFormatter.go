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
	debug              bool
}

func NewDefaultFormatterWithParameters(printer interfaces.SqlPrinter, parameterLookup map[int]string) *DefaultFormatter {
	return &DefaultFormatter{
		printer:            printer,
		detectedParameters: parameterLookup,
		debug:              false,
	}
}

func NewDefaultFormatter(printer interfaces.SqlPrinter) *DefaultFormatter {
	return NewDefaultFormatterWithParameters(printer, nil)
}

func (df *DefaultFormatter) d() {
	if !df.debug {
		return
	}

	fmt.Println(df.printer)
}

func (df *DefaultFormatter) p(msg string) {
	// Dump the printer's content before we panic to aid with debugging
	fmt.Println(df.printer)

	panic(msg + " not supported")
}

func (df *DefaultFormatter) String() string {
	return df.printer.String()
}

func (df *DefaultFormatter) PrintWithClause(wc nodes.WithClause) {
	if wc.Recursive {
		df.p("With Clause - Recursive")
	}

	for i := range wc.Ctes.Items {
		if i == 0 {
			df.printer.PrintKeyword("with ")
		}

		df.printNode(wc.Ctes.Items[i], false)

		if i < len(wc.Ctes.Items)-1 {
			df.printer.PrintString(",")
		}

		df.printer.NewLine()
	}
}

func (df *DefaultFormatter) PrintSelectStatementTargets(ss nodes.SelectStmt) {
	printDistinct := false

	for i := range ss.DistinctClause.Items {
		if ss.DistinctClause.Items[i] == nil {
			df.printer.PrintKeyword("distinct ", true)
			printDistinct = true
		} else {
			df.printNode(ss.DistinctClause.Items[i], !printDistinct)
			printDistinct = false
		}
	}

	for i := range ss.TargetList.Items {
		df.printNode(ss.TargetList.Items[i], !printDistinct)
		printDistinct = false

		if i < len(ss.TargetList.Items)-1 {
			df.printer.PrintString(",")
			df.printer.NewLine()
		}
	}

	if ss.IntoClause != nil {
		df.p("Select - Into clause")
	}

	df.printer.DecIndent()

	// only drop a new line if we are selecting something
	if len(ss.DistinctClause.Items) > 0 || len(ss.TargetList.Items) > 0 {
		df.printer.NewLine()
	}
}

func (df *DefaultFormatter) PrintSelectStatementFromClause(ss nodes.SelectStmt) {
	for i := range ss.FromClause.Items {
		if je, ok := ss.FromClause.Items[i].(nodes.JoinExpr); ok {
			df.PrintJoin(i == 0, je)
		} else {
			if i == 0 {
				df.printer.PrintKeyword("from", true)
				df.printer.IncIndent()
			}

			df.printer.NewLine()
			df.printNode(ss.FromClause.Items[i], true)

			if i < len(ss.FromClause.Items)-1 {
				df.printer.PrintString(",")
			}
		}
	}
}

func (df *DefaultFormatter) PrintSelectStatementWhereClause(ss nodes.SelectStmt) {
	df.printer.DecIndent()
	if ss.WhereClause != nil {
		df.printer.NewLine()
		df.printer.PrintKeyword("where", true)
		df.printer.NewLine()
		df.printer.IncIndent()
		df.printNode(ss.WhereClause, true)
		df.printer.DecIndent()
	}
}

func (df *DefaultFormatter) PrintSelectStatementSortClause(ss nodes.SelectStmt) {
	if len(ss.SortClause.Items) > 0 {
		df.printer.NewLine()
		df.printer.PrintKeyword("order by", true)
		df.printer.NewLine()
		df.printer.IncIndent()

		for i, item := range ss.SortClause.Items {
			df.printNode(item, true)
			if i < len(ss.SortClause.Items)-1 {
				df.printer.PrintString(",")
				df.printer.NewLine()
			}
		}
		df.printer.DecIndent()
	}
}

func (df *DefaultFormatter) PrintSelectStatementLimitClause(ss nodes.SelectStmt) {
	if ss.LimitCount != nil {
		df.printer.NewLine()
		df.printer.PrintKeyword("limit", true)
		df.printer.NewLine()
		df.printer.IncIndent()
		df.printNode(ss.LimitCount, true)
		df.printer.DecIndent()
	}
}

func (df *DefaultFormatter) PrintSelectStatementGroupByClause(ss nodes.SelectStmt) {
	if len(ss.GroupClause.Items) > 0 {
		df.printer.NewLine()
		df.printer.PrintKeyword("group by", true)
		df.printer.IncIndent()
		df.printer.NewLine()

		for i, item := range ss.GroupClause.Items {
			df.printNode(item, true)
			if i < len(ss.GroupClause.Items)-1 {
				df.printer.PrintString(",")
				df.printer.NewLine()
			}
		}

		df.printer.DecIndent()
	}
}

func (df *DefaultFormatter) PrintSelectStatementHavingClause(ss nodes.SelectStmt) {
	if ss.HavingClause != nil {
		df.printer.NewLine()
		df.printer.PrintKeyword("having", true)
		df.printer.IncIndent()
		df.printer.NewLine()
		df.printNode(ss.HavingClause, true)
	}
}

func (df *DefaultFormatter) PrintSelectStatement(ss nodes.SelectStmt) {
	if ss.WithClause != nil {
		df.PrintWithClause(*ss.WithClause)
	}

	df.printer.PrintKeyword("select", true)
	df.printer.NewLine()
	df.printer.IncIndent()

	df.PrintSelectStatementTargets(ss)
	df.PrintSelectStatementFromClause(ss)
	df.PrintSelectStatementWhereClause(ss)
	df.PrintSelectStatementGroupByClause(ss)
	df.PrintSelectStatementHavingClause(ss)
	df.PrintSelectStatementSortClause(ss)
	df.PrintSelectStatementLimitClause(ss)
}

func (df *DefaultFormatter) PrintResTarget(nt nodes.ResTarget, withIndent bool) {
	retVal := ""

	if nt.Name != nil {
		retVal = *nt.Name
	}

	for i := range nt.Indirection.Items {
		df.printNode(nt.Indirection.Items[i], withIndent)
	}

	df.printNode(nt.Val, withIndent)

	if len(retVal) > 0 {
		df.printer.PrintKeyword(" as ")
		df.printer.PrintString("\"" + retVal + "\"")
	}
}

func (df *DefaultFormatter) PrintColumnRef(cr nodes.ColumnRef, withIndent bool) {
	for i := range cr.Fields.Items {
		df.printNode(cr.Fields.Items[i], withIndent && i == 0)

		if i < len(cr.Fields.Items)-1 {
			df.printer.PrintString(".")
		}
	}
}

func (df *DefaultFormatter) PrintJoin(first bool, join nodes.JoinExpr) {
	if first {
		df.printer.PrintKeyword("from", true)
		df.printer.NewLine()
		df.printer.IncIndent()
	}

	if join.IsNatural {
		df.p("Join - Natural")
	}

	// cross join
	if join.Jointype == nodes.JOIN_INNER && join.Quals == nil {
		df.printNode(join.Larg, true)
		df.printer.NewLine()
		df.printer.DecIndent()
		df.printer.PrintKeyword("cross join", true)
		df.printer.NewLine()
		df.printer.IncIndent()
		df.printNode(join.Rarg, true)
	} else {
		df.printNode(join.Larg, true)
		df.printer.NewLine()
		df.printer.DecIndent()
		df.PrintJoinType(join.Jointype, true)

		newLine := true

		switch join.Rarg.(type) {
		case nodes.RangeSubselect:
			newLine = false
		}

		// Sub-select's manage their own newlines, due to lateral joins
		if newLine {
			df.printer.NewLine()
			df.printer.IncIndent()
		}

		df.printNode(join.Rarg, true)

		if len(join.UsingClause.Items) > 0 {
			df.p("Join - Using Clause")
		}

		df.printer.NewLine()
		df.printer.PrintKeyword("on", true)
		df.printer.NewLine()
		df.printer.IncIndent()
		df.printNode(join.Quals, true)
		df.printer.DecIndent()
	}
}

func (df *DefaultFormatter) PrintAlias(alias nodes.Alias) {
	if alias.Aliasname != nil {
		df.printer.PrintString(*(alias.Aliasname))
	}

	if len(alias.Colnames.Items) > 0 {
		df.printer.PrintString("(")

		for i, col := range alias.Colnames.Items {
			df.printNode(col, false)
			if i < len(alias.Colnames.Items)-1 {
				df.printer.PrintString(", ")
			}
		}

		df.printer.PrintString(")")
	}
}

func (df *DefaultFormatter) PrintJoinType(joinType nodes.JoinType, withIndent bool) {
	jt := ""

	switch joinType {
	case nodes.JOIN_INNER:
		jt = "join"

	case nodes.JOIN_LEFT:
		jt = "left join"

	case nodes.JOIN_FULL:
		jt = "full join"

	case nodes.JOIN_RIGHT:
		jt = "right join"

	default:
		df.p(fmt.Sprintf("join type - %+v not supported", joinType))
	}

	if len(jt) > 0 {
		df.printer.PrintKeyword(jt, withIndent)
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

	df.printer.PrintString(name, withIndent)

	if rv.Alias != nil {
		df.printer.PrintString(" ")
		df.PrintAlias(*rv.Alias)
	}
}

func (df *DefaultFormatter) PrintBoolExpr(be nodes.BoolExpr, prevOp nodes.BoolExprType, withIndent bool) {
	if be.Xpr != nil {
		df.printNode(be.Xpr, withIndent)
		df.printer.PrintString(" ")
		df.PrintBoolExprType(be.Boolop, false)
	}

	// Drop parentheses when change operators
	parentheses := be.Boolop != prevOp && be.Boolop != 2

	if parentheses {
		df.printer.PrintString("(", withIndent)

		df.printer.NewLine()
		df.printer.IncIndent()
	}

	for i := range be.Args.Items {
		if int(be.Boolop) == 2 {
			df.PrintBoolExprType(be.Boolop, withIndent)
		}

		if tbe, ok := be.Args.Items[i].(nodes.BoolExpr); ok {
			df.PrintBoolExpr(tbe, be.Boolop, withIndent || i == 0)
		} else {
			// only print an indent if this is the first item and the operator is not a not
			df.printNode(be.Args.Items[i], i == 0 && (int(be.Boolop) != 2))
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
		df.printer.PrintString(")", true)
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
		df.printer.PrintKeyword(kw, withIndent)
	}
}

func (df *DefaultFormatter) PrintAExpr(ae nodes.A_Expr, withIndent bool) {
	if ae.Lexpr != nil {
		df.printNode(ae.Lexpr, withIndent)
	}

	switch ae.Kind {
	case nodes.AEXPR_OP_ANY:
		df.PrintAExprKeywords(ae.Name, true)
		df.printer.PrintKeyword("any")
		df.printer.PrintString("(")
		df.printNode(ae.Rexpr, false)
		df.printer.PrintString(")")

	case nodes.AEXPR_OP_ALL:
		df.PrintAExprKeywords(ae.Name, true)
		df.printer.PrintKeyword("all")
		df.printer.PrintString("(")
		df.printNode(ae.Rexpr, false)
		df.printer.PrintString(")")

	case nodes.AEXPR_BETWEEN,
		nodes.AEXPR_NOT_BETWEEN:

		df.PrintAExprKeywords(ae.Name, true)
		df.printNode(ae.Rexpr.(nodes.List).Items[0], false)
		df.printer.PrintKeyword(" and ")
		df.printNode(ae.Rexpr.(nodes.List).Items[1], false)

	case nodes.AEXPR_LIKE:
		df.printer.PrintKeyword(" like ")
		df.printNode(ae.Rexpr, false)

	case nodes.AEXPR_ILIKE:
		df.printer.PrintKeyword(" ilike ")
		df.printNode(ae.Rexpr, false)

	case nodes.AEXPR_OP:
		if ae.Lexpr == nil {
			// we need an indent, that Lexpr gives
			df.printer.PrintString("", true)
		}
		// Prints the operator
		df.PrintAExprKeywords(ae.Name, ae.Lexpr != nil)

		df.printNode(ae.Rexpr, false)

	default:
		df.p(fmt.Sprintf("AExpr: %v", ae.Kind))
	}
}

func (df *DefaultFormatter) PrintAExprKeywords(op nodes.List, spaces bool) {
	for i := range op.Items {
		// These are keywords
		if spaces {
			df.printer.PrintString(" ")
		}

		df.printer.PrintKeyword(op.Items[i].(nodes.String).Str)

		if spaces {
			df.printer.PrintString(" ")
		}
	}
}

func (df *DefaultFormatter) PrintNullTest(nt nodes.NullTest, withIndent bool) {
	if nt.Xpr != nil {
		df.printNode(nt.Xpr, withIndent)
	}

	df.printNode(nt.Arg, withIndent)

	if nt.Nulltesttype == nodes.IS_NULL {
		df.printer.PrintKeyword(" is null")
	} else {
		df.printer.PrintKeyword(" is not null")
	}
}

func (df *DefaultFormatter) PrintAConst(ac nodes.A_Const, withindent bool) {
	switch ac.Val.(type) {
	case nodes.String:

		df.printer.PrintString("'", withindent)
		df.printNode(ac.Val, false)
		df.printer.PrintString("'")

	default:

		df.printNode(ac.Val, withindent)
	}
}

func (df *DefaultFormatter) PrintCommonTableExpr(cte nodes.CommonTableExpr) {
	if cte.Ctename != nil {
		df.printer.PrintString(*cte.Ctename)
		df.printer.PrintKeyword(" as ")
		df.printer.PrintString("(")
		df.printer.NewLine()
		df.printer.IncIndent()
	}

	if cte.Cterecursive {
		df.p("CTE - Recursive")
	}

	df.printNode(cte.Ctequery, false)

	df.printer.NewLine()
	df.printer.DecIndent()
	df.printer.PrintString(")", true)
}

func (df *DefaultFormatter) PrintParamRef(pr nodes.ParamRef, withIndent bool) {
	if df.detectedParameters != nil {
		if param, ok := df.detectedParameters[df.paramCounter]; ok {
			df.printer.PrintString(param, withIndent)
			df.paramCounter++
			return
		}
	}

	df.printer.PrintString("?", withIndent)
}

func (df *DefaultFormatter) PrintSubSelect(ss nodes.RangeSubselect, withIndent bool) {
	if ss.Lateral {
		df.printer.PrintKeyword(" lateral")
	}
	df.printer.NewLine()
	df.printer.IncIndent()
	df.printer.PrintString("(", true)
	df.printer.NewLine()
	df.printer.IncIndent()
	df.printNode(ss.Subquery, withIndent)
	df.printer.NewLine()
	df.printer.DecIndent()
	df.printer.PrintString(") ", true)

	if ss.Alias != nil {
		df.PrintAlias(*ss.Alias)
	}
}

func (df *DefaultFormatter) PrintTypeCast(tc nodes.TypeCast, withIndent bool) {
	// all of this garbage is required to convert an optimized 't'::bool back to a true
	// because only crazy people prefer the later.
	isTrue := false
	isAConst := false

	switch tc.Arg.(type) {
	case nodes.A_Const:
		isAConst = true
	}

	if isAConst {
		isString := false
		ac := tc.Arg.(nodes.A_Const)
		switch ac.Val.(type) {
		case nodes.String:
			isString = true
		}

		if isString && ac.Val.(nodes.String).Str == "t" {
			isTrue = true
		}
	}

	if isTrue {
		df.printer.PrintString("true", withIndent)
		return
	}

	df.printNode(tc.Arg, withIndent)
	if tc.TypeName != nil {
		df.printer.PrintString("::")
		df.PrintTypeName(*tc.TypeName)
	}
}

func (df *DefaultFormatter) PrintTypeName(tn nodes.TypeName) {
	if len(tn.Names.Items) > 1 {
		df.printNode(tn.Names.Items[1], false)
	}
}

func (df *DefaultFormatter) PrintSubLink(sl nodes.SubLink, withIndent bool) {
	switch sl.SubLinkType {
	case nodes.ROWCOMPARE_SUBLINK,
		nodes.MULTIEXPR_SUBLINK,
		nodes.ARRAY_SUBLINK:

		df.p(fmt.Sprintf("Unsupported sublinktype: %+v", sl.SubLinkType))

	case nodes.ANY_SUBLINK:

		if sl.Testexpr != nil {
			df.printNode(sl.Testexpr, withIndent)
		}

		if len(sl.OperName.Items) == 0 {
			df.printer.PrintKeyword(" in")
		} else {
			df.printer.PrintString(" ")

			for _, opp := range sl.OperName.Items {
				df.printNode(opp, false)
			}

			df.printer.PrintString(" ")
			df.printer.PrintKeyword("any")
		}

	case nodes.ALL_SUBLINK:

		if sl.Testexpr != nil {
			df.printNode(sl.Testexpr, withIndent)
		}

		df.printer.PrintString(" ")

		for _, opp := range sl.OperName.Items {
			df.printNode(opp, false)
		}

		df.printer.PrintString(" ")
		df.printer.PrintKeyword("all")

	case nodes.EXISTS_SUBLINK:

		df.printer.PrintKeyword("exists", withIndent)

	case nodes.EXPR_SUBLINK:

		df.printer.PrintString("", true)
	}

	df.printer.PrintString("(")
	df.printer.NewLine()
	df.printer.IncIndent()
	df.printNode(sl.Subselect, withIndent)
	df.printer.NewLine()
	df.printer.DecIndent()
	df.printer.PrintString(")", true)
}

func (df *DefaultFormatter) PrintRangeFunction(rf nodes.RangeFunction, withIndent bool) {
	for _, f := range rf.Functions.Items {
		df.printNode(f, withIndent)
	}

	if rf.Alias != nil {
		df.printer.PrintString(" ")
		df.printNode(*rf.Alias, false)
	}
}

func (df *DefaultFormatter) PrintFuncCallName(fc nodes.FuncCall, withIndent bool) {
	for _, name := range fc.Funcname.Items {
		if s, ok := name.(nodes.String); ok {
			df.printer.PrintFunction(s.Str, withIndent)
		} else {
			df.printNode(name, withIndent)
		}
	}
}

func (df *DefaultFormatter) PrintFuncCallArgs(fc nodes.FuncCall, withIndent bool) {
	for i, arg := range fc.Args.Items {
		df.printNode(arg, false)
		if i < len(fc.Args.Items)-1 {
			df.printer.PrintString(", ")
		}
	}
}

func (df *DefaultFormatter) PrintFuncCallOrder(fc nodes.FuncCall, withIndent bool) {
	if len(fc.AggOrder.Items) > 0 {
		df.printer.PrintKeyword(" order by")
		for i, item := range fc.AggOrder.Items {
			df.printNode(item, false)
			if i < len(fc.AggOrder.Items)-1 {
				df.printer.PrintString(", ")
			}
		}
	}
}

func (df *DefaultFormatter) PrintFuncCallAggFilter(fc nodes.FuncCall, withIndent bool) {
	if fc.AggFilter != nil {
		df.printer.NewLine()
		df.printer.IncIndent()
		df.printer.PrintKeyword("filter ", true)
		df.printer.PrintString("(")
		df.printer.NewLine()
		df.printer.IncIndent()
		df.printer.PrintKeyword("where", true)
		df.printer.NewLine()
		df.printer.IncIndent()
		df.printNode(fc.AggFilter, true)
		df.printer.DecIndent()
		df.printer.DecIndent()
		df.printer.NewLine()
		df.printer.PrintString(")", true)
		df.printer.DecIndent()
	}
}

func (df *DefaultFormatter) PrintFuncCall(fc nodes.FuncCall, withIndent bool) {
	df.PrintFuncCallName(fc, withIndent)

	df.printer.PrintString("(")

	if fc.AggDistinct {
		df.printer.PrintKeyword("distinct ")
	}

	if fc.AggStar {
		df.printer.PrintString("*")
	}

	df.PrintFuncCallArgs(fc, withIndent)
	df.PrintFuncCallOrder(fc, withIndent)
	df.printer.PrintString(")")
	df.PrintFuncCallAggFilter(fc, withIndent)
}

func (df *DefaultFormatter) PrintSortBy(sb nodes.SortBy, withIndent bool) {
	df.printNode(sb.Node, withIndent)
	if sb.SortbyDir == nodes.SORTBY_DESC {
		df.printer.PrintKeyword(" desc")
	}
	if sb.SortbyNulls == nodes.SORTBY_NULLS_LAST {
		df.printer.PrintKeyword(" nulls last")
	}
}

// PrintNode This is the main entry point for the AST crawler
func (df *DefaultFormatter) PrintNode(node nodes.Node) {
	df.printNode(node, false)
}

// printNode This is the main forking function that decides what to do with a node.
//
// Note: Other `node` printers call this function to print `node` objects, so this is a generic `node` printer
func (df *DefaultFormatter) printNode(node nodes.Node, withIndent bool) {
	if node == nil {
		return
	}

	switch node.(type) {
	case nodes.RawStmt:
		df.printNode(node.(nodes.RawStmt).Stmt, withIndent)

	case nodes.SelectStmt:
		df.PrintSelectStatement(node.(nodes.SelectStmt))

	case nodes.IntoClause:
		df.printer.PrintKeyword(" into ")

	case nodes.ResTarget:
		df.PrintResTarget(node.(nodes.ResTarget), withIndent)

	case nodes.ColumnRef:
		df.PrintColumnRef(node.(nodes.ColumnRef), withIndent)

	case nodes.A_Star:
		df.printer.PrintString("*", withIndent)

	case nodes.String:
		df.printer.PrintString(node.(nodes.String).Str, withIndent)

	case nodes.JoinExpr:
		df.PrintJoin(false, node.(nodes.JoinExpr))

	case nodes.RangeVar:
		df.PrintRangeVar(node.(nodes.RangeVar), withIndent)

	case nodes.BoolExpr:
		df.PrintBoolExpr(node.(nodes.BoolExpr), node.(nodes.BoolExpr).Boolop, withIndent)

	case nodes.A_Expr:
		df.PrintAExpr(node.(nodes.A_Expr), withIndent)

	case nodes.A_Const:
		df.PrintAConst(node.(nodes.A_Const), withIndent)

	case nodes.Integer:
		df.printer.PrintInt64(node.(nodes.Integer).Ival, withIndent)

	case nodes.Float:
		val, err := strconv.ParseFloat(node.(nodes.Float).Str, 64)
		if err != nil {
			df.p(node.(nodes.Float).Str + " failed to parse as a float")
		}

		df.printer.PrintFloat64(val, withIndent)

	case nodes.NullTest:
		df.PrintNullTest(node.(nodes.NullTest), withIndent)

	case nodes.CommonTableExpr:
		df.PrintCommonTableExpr(node.(nodes.CommonTableExpr))

	case nodes.ParamRef:
		df.PrintParamRef(node.(nodes.ParamRef), withIndent)

	case nodes.RangeSubselect:
		df.PrintSubSelect(node.(nodes.RangeSubselect), withIndent)

	case nodes.SubLink:
		df.PrintSubLink(node.(nodes.SubLink), withIndent)

	case nodes.Alias:
		df.PrintAlias(node.(nodes.Alias))

	case nodes.TypeCast:
		df.PrintTypeCast(node.(nodes.TypeCast), withIndent)

	case nodes.RangeFunction:
		df.PrintRangeFunction(node.(nodes.RangeFunction), withIndent)

	case nodes.List:
		list := node.(nodes.List)
		for _, item := range list.Items {
			df.printNode(item, withIndent)
		}

	case nodes.FuncCall:
		df.PrintFuncCall(node.(nodes.FuncCall), withIndent)

	case nodes.SortBy:
		df.PrintSortBy(node.(nodes.SortBy), withIndent)

	default:
		df.p(fmt.Sprintf("Node: %T", node))
	}
}
