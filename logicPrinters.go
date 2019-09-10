package main

import (
	"fmt"
	"regexp"

	nodes "github.com/lfittl/pg_query_go/nodes"
)

func printWithClause(wc nodes.WithClause) {
	if wc.Recursive {
		p()
	}

	for i := range wc.Ctes.Items {
		if i == 0 {
			formatter.PrintKeywordNoIndent("with ")
		}

		printNode(wc.Ctes.Items[i], false)

		if i < len(wc.Ctes.Items)-1 {
			formatter.PrintString(", ")
		}
	}
}

func printSelectStatement(ss nodes.SelectStmt) {
	if ss.WithClause != nil {
		printWithClause(*ss.WithClause)
	}

	formatter.PrintKeyword("select")
	formatter.NewLine()
	formatter.IncIndent()

	for i := range ss.DistinctClause.Items {
		if ss.DistinctClause.Items[i] == nil {
			formatter.PrintKeywordNoIndent("distinct")
			formatter.NewLine()
		} else {
			printNode(ss.DistinctClause.Items[i], true)
		}
	}

	if ss.IntoClause != nil {
		p()
	}

	for i := range ss.TargetList.Items {
		printNode(ss.TargetList.Items[i], true)

		if i < len(ss.TargetList.Items)-1 {
			formatter.PrintStringNoIndent(",")
			formatter.NewLine()
		}
	}

	formatter.DecIndent()
	formatter.NewLine()

	for i := range ss.FromClause.Items {
		if je, ok := ss.FromClause.Items[i].(nodes.JoinExpr); ok {
			printJoin(i == 0, je)
		} else {
			formatter.PrintKeyword("from")
			formatter.NewLine()
			formatter.IncIndent()
			printNode(ss.FromClause.Items[i], true)
		}

		formatter.NewLine()
	}

	formatter.DecIndent()
	formatter.PrintKeyword("where")
	formatter.NewLine()
	formatter.IncIndent()
	printNode(ss.WhereClause, true)
	formatter.DecIndent()
}

func nilCheck(n nodes.Node) bool {
	return n == nil
}

func printResTarget(nt nodes.ResTarget, withIndent bool) {
	retVal := ""

	if nt.Name != nil {
		retVal = *nt.Name
	}

	for i := range nt.Indirection.Items {
		printNode(nt.Indirection.Items[i], withIndent)
	}

	printNode(nt.Val, withIndent)

	formatter.PrintStringNoIndent(retVal)
}

func printColumnRef(cr nodes.ColumnRef, withIndent bool) {
	for i := range cr.Fields.Items {
		printNode(cr.Fields.Items[i], withIndent)

		if i < len(cr.Fields.Items)-1 {
			formatter.PrintStringNoIndent(".")
		}

		// negate indentation on consecutive runs to stop tab.<space>col prints
		withIndent = false
	}
}

func printJoin(first bool, join nodes.JoinExpr) {
	if first {
		formatter.PrintKeyword("from")
		formatter.NewLine()
		formatter.IncIndent()
	} else {
		formatter.IncIndent()
		printJoinType(join.Jointype, true)
	}

	if join.IsNatural {
		p()
	}

	// cross join
	if join.Jointype == nodes.JOIN_INNER && join.Quals == nil {
		printNode(join.Larg, true)
		formatter.NewLine()
		formatter.DecIndent()
		formatter.PrintKeyword("cross join ")
		printNode(join.Rarg, false)
		formatter.IncIndent()
	} else {
		printNode(join.Larg, true)
		formatter.NewLine()
		formatter.DecIndent()
		printJoinType(join.Jointype, true)
		printNode(join.Rarg, false)

		if len(join.UsingClause.Items) > 0 {
			p()
		}

		formatter.NewLine()
		formatter.IncIndent()
		formatter.PrintKeyword("on ")
		printNode(join.Quals, false)
	}
}

func printAlias(alias *nodes.Alias) {
	if alias == nil {
		return
	}

	if alias.Aliasname != nil {
		formatter.PrintStringNoIndent(*(alias.Aliasname))
	}

	if len(alias.Colnames.Items) > 0 {
		p()
	}
}

func printJoinType(joinType nodes.JoinType, withIndent bool) {
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
		p()

	case nodes.JOIN_ANTI:
		p()

	case nodes.JOIN_UNIQUE_OUTER:
		p()

	case nodes.JOIN_UNIQUE_INNER:
		p()
	}

	if len(jt) > 0 {
		if withIndent {
			formatter.PrintKeyword(jt)
		} else {
			formatter.PrintKeywordNoIndent(jt)
		}
	}
}

func printRangeVar(rv nodes.RangeVar, withIndent bool) {
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
		formatter.PrintString(name)
	} else {
		formatter.PrintStringNoIndent(name)
	}

	if rv.Alias != nil {
		formatter.PrintStringNoIndent(" ")
		printAlias(rv.Alias)
	}
}

func printBoolExpr(be nodes.BoolExpr, prevOp nodes.BoolExprType, withIndent bool) {
	if be.Xpr != nil {
		printNode(be.Xpr, withIndent)
		formatter.PrintStringNoIndent(" ")
		printBoolExprType(be.Boolop, false)
	}

	// Drop parentheses when change operators
	parentheses := be.Boolop != prevOp && be.Boolop != 2

	if parentheses {
		if withIndent {
			formatter.PrintString("(")
		} else {
			formatter.PrintStringNoIndent("(")
		}

		formatter.NewLine()
		formatter.IncIndent()
	}

	for i := range be.Args.Items {
		if int(be.Boolop) == 2 {
			printBoolExprType(be.Boolop, withIndent)
		}

		if tbe, ok := be.Args.Items[i].(nodes.BoolExpr); ok {
			printBoolExpr(tbe, be.Boolop, withIndent)
		} else {
			// only print an indent if this is the first item and the operator is not a not
			printNode(be.Args.Items[i], i == 0 && (int(be.Boolop) != 2))
		}

		if i < len(be.Args.Items)-1 {
			formatter.NewLine()
			printBoolExprType(be.Boolop, true)
		}

		// Stop printing consecutive args with indent to stop crap like AND<space>tab.col...
		withIndent = false
	}

	if parentheses {
		formatter.NewLine()
		formatter.DecIndent()
		formatter.PrintString(")")
	}
}

func printBoolExprType(exprType nodes.BoolExprType, withIndent bool) {
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
			formatter.PrintKeyword(kw)
		} else {
			formatter.PrintKeywordNoIndent(kw)
		}
	}
}

func printAExpr(ae nodes.A_Expr, withIndent bool) {
	printNode(ae.Lexpr, withIndent)

	switch ae.Kind {
	case nodes.AEXPR_OP_ANY:
		printAExprKeywords(ae.Name)
		formatter.PrintKeywordNoIndent("any")
		formatter.PrintStringNoIndent("( ")
		printNode(ae.Rexpr, false)
		formatter.PrintStringNoIndent(" )")

	case nodes.AEXPR_OP_ALL:
		printAExprKeywords(ae.Name)
		formatter.PrintKeywordNoIndent("all")
		formatter.PrintStringNoIndent("( ")
		printNode(ae.Rexpr, false)
		formatter.PrintStringNoIndent(" )")

	case nodes.AEXPR_BETWEEN,
		nodes.AEXPR_NOT_BETWEEN:

		printAExprKeywords(ae.Name)
		printNode(ae.Rexpr.(nodes.List).Items[0], false)
		formatter.PrintKeywordNoIndent(" and ")
		printNode(ae.Rexpr.(nodes.List).Items[1], false)

	case nodes.AEXPR_LIKE:
		formatter.PrintKeywordNoIndent(" like ")
		printNode(ae.Rexpr, false)

	case nodes.AEXPR_ILIKE:
		formatter.PrintKeywordNoIndent(" ilike ")
		printNode(ae.Rexpr, false)

	case nodes.AEXPR_OP:
		// Prints the operator
		printAExprKeywords(ae.Name)

		printNode(ae.Rexpr, false)

	default:
		fmt.Printf("\nAExpr: %v\n", ae.Kind)
		p()

	}
}

func printAExprKeywords(op nodes.List) {
	for i := range op.Items {
		// These are keywords
		formatter.PrintStringNoIndent(" ")
		formatter.PrintKeywordNoIndent(op.Items[i].(nodes.String).Str)
		formatter.PrintStringNoIndent(" ")
	}
}

func printNullTest(nt nodes.NullTest, withIndent bool) {
	if nt.Xpr != nil {
		printNode(nt.Xpr, withIndent)
	}

	printNode(nt.Arg, withIndent)

	if nt.Nulltesttype == nodes.IS_NULL {
		formatter.PrintKeywordNoIndent(" is null")
	} else {
		formatter.PrintKeywordNoIndent(" is not null")
	}
}

func printAConst(ac nodes.A_Const) {
	switch ac.Val.(type) {
	case nodes.String:
		formatter.PrintStringNoIndent("'")
		printNode(ac.Val, false)
		formatter.PrintStringNoIndent("'")

	default:
		printNode(ac.Val, false)
	}
}

func printCommonTableExpr(cte nodes.CommonTableExpr) {
	if cte.Ctename != nil {
		formatter.PrintStringNoIndent(*cte.Ctename)
		formatter.PrintStringNoIndent(" (")
		formatter.NewLine()
		formatter.IncIndent()
	}

	if cte.Cterecursive {
		p()
	}

	printNode(cte.Ctequery, false)

	formatter.NewLine()
	formatter.DecIndent()
	formatter.PrintString(")")
	formatter.NewLine()
}

func printParamRef(pr nodes.ParamRef, withIndent bool) {
	if param, ok := detectedParameters[paramCounter]; ok {
		if withIndent {
			formatter.PrintString(param)
		} else {
			formatter.PrintStringNoIndent(param)
		}
		paramCounter++
	} else {
		if withIndent {
			formatter.PrintString("?")
		} else {
			formatter.PrintStringNoIndent("?")
		}
	}
}

func printNode(node nodes.Node, withIndent bool) {
	switch node.(type) {
	case nodes.RawStmt:
		printNode(node.(nodes.RawStmt).Stmt, withIndent)

	case nodes.SelectStmt:
		printSelectStatement(node.(nodes.SelectStmt))

	case nodes.IntoClause:
		formatter.PrintKeywordNoIndent(" into ")

	case nodes.ResTarget:
		printResTarget(node.(nodes.ResTarget), withIndent)

	case nodes.ColumnRef:
		printColumnRef(node.(nodes.ColumnRef), withIndent)

	case nodes.A_Star:
		formatter.PrintString("*")

	case nodes.String:
		if withIndent {
			formatter.PrintString(node.(nodes.String).Str)
		} else {
			formatter.PrintStringNoIndent(node.(nodes.String).Str)
		}

	case nodes.JoinExpr:
		printJoin(false, node.(nodes.JoinExpr))

	case nodes.RangeVar:
		printRangeVar(node.(nodes.RangeVar), withIndent)

	case nodes.BoolExpr:
		printBoolExpr(node.(nodes.BoolExpr), node.(nodes.BoolExpr).Boolop, withIndent)

	case nodes.A_Expr:
		printAExpr(node.(nodes.A_Expr), withIndent)

	case nodes.A_Const:
		printAConst(node.(nodes.A_Const))

	case nodes.Integer:
		if withIndent {
			formatter.PrintInt(node.(nodes.Integer).Ival)
		} else {
			formatter.PrintIntNoIndent(node.(nodes.Integer).Ival)
		}

	case nodes.NullTest:
		printNullTest(node.(nodes.NullTest), withIndent)

	case nodes.CommonTableExpr:
		printCommonTableExpr(node.(nodes.CommonTableExpr))

	case nodes.ParamRef:
		printParamRef(node.(nodes.ParamRef), withIndent)

	default:
		fmt.Printf("NOT SUPPORTED : %T\n", node)
		p()
	}
}

func processNamedParameters(sql string) string {
	compRegEx := regexp.MustCompile(`(\?\w+)`)

	for i, match := range compRegEx.FindAllString(sql, -1) {
		detectedParameters[i] = match
	}

	return compRegEx.ReplaceAllString(sql, "?")
}
