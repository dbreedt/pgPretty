package main

import (
	"fmt"
	"regexp"
	"strings"

	pg_query "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
)

type caseFormatter func(s string) string

var indent = "  "
var keywordFormatter caseFormatter = strings.ToLower
var detectedParameters = make(map[int]string)
var paramCounter int

func p() {
	panic("not supported")
}

func printIndent(num int) {
	fmt.Printf(strings.Repeat(indent, num))
}

func printWithClause(wc nodes.WithClause) {
	if wc.Recursive {
		p()
	}

	for i := range wc.Ctes.Items {
		if i == 0 {
			fmt.Printf(keywordFormatter("with "))
		}

		printNode(wc.Ctes.Items[i])

		if i < len(wc.Ctes.Items)-1 {
			fmt.Printf(", ")
		}
	}
}

func printSelectStatement(ss nodes.SelectStmt) string {
	if ss.WithClause != nil {
		printWithClause(*ss.WithClause)
	}

	fmt.Printf("%s\n", keywordFormatter("select"))

	for i := range ss.DistinctClause.Items {
		if ss.DistinctClause.Items[i] == nil {
			printIndent(1)
			fmt.Printf("%s\n", keywordFormatter("distinct"))
		} else {
			printNode(ss.DistinctClause.Items[i])
		}
	}

	if ss.IntoClause != nil {
		p()
	}

	for i := range ss.TargetList.Items {
		printIndent(1)
		printNode(ss.TargetList.Items[i])

		if i < len(ss.TargetList.Items)-1 {
			fmt.Printf(",\n")
		}
	}

	fmt.Println()

	for i := range ss.FromClause.Items {
		if je, ok := ss.FromClause.Items[i].(nodes.JoinExpr); ok {
			printJoin(i == 0, je)
		} else {
			fmt.Printf("%s\n", keywordFormatter("from"))
			printIndent(1)
			printNode(ss.FromClause.Items[i])
		}
		fmt.Println()
	}

	fmt.Printf("%s\n", keywordFormatter("where"))
	printIndent(1)
	printNode(ss.WhereClause)

	return ""
}

func nilCheck(n nodes.Node) bool {
	return n == nil
}

func printResTarget(nt nodes.ResTarget) string {
	retVal := ""

	if nt.Name != nil {
		retVal = *nt.Name
	}

	for i := range nt.Indirection.Items {
		printNode(nt.Indirection.Items[i])
	}

	printNode(nt.Val)

	return retVal
}

func printColumRef(cr nodes.ColumnRef) string {
	for i := range cr.Fields.Items {
		printNode(cr.Fields.Items[i])

		if i < len(cr.Fields.Items)-1 {
			fmt.Printf(".")
		}
	}

	return ""
}

func printJoin(first bool, join nodes.JoinExpr) {
	if first {
		fmt.Printf("%s\n", keywordFormatter("from"))
		printIndent(1)
	} else {
		printIndent(1)
		printJoinType(join.Jointype)
	}

	if join.IsNatural {
		p()
	}

	// cross join
	if join.Jointype == nodes.JOIN_INNER && join.Quals == nil {
		printNode(join.Larg)
		fmt.Println()
		printIndent(1)
		fmt.Printf("%s ", keywordFormatter("cross join"))
		printNode(join.Rarg)
	} else {

		printNode(join.Larg)
		println()
		printJoinType(join.Jointype)
		printNode(join.Rarg)

		if len(join.UsingClause.Items) > 0 {
			p()
		}

		fmt.Println()
		printIndent(1)
		fmt.Printf(keywordFormatter("on "))
		printNode(join.Quals)
	}
}

func printAlias(alias *nodes.Alias) {
	if alias == nil {
		return
	}

	if alias.Aliasname != nil {
		fmt.Printf(*(alias.Aliasname))
	}

	if len(alias.Colnames.Items) > 0 {
		p()
	}
}

func printJoinType(joinType nodes.JoinType) {
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

	fmt.Printf(keywordFormatter(jt))
}

func printRangeVar(rv nodes.RangeVar) {
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

	fmt.Printf(name)

	if rv.Alias != nil {
		fmt.Printf(" ")
		printAlias(rv.Alias)
	}
}

func printBoolExpr(be nodes.BoolExpr, prevOp nodes.BoolExprType) {
	if be.Xpr != nil {
		printNode(be.Xpr)
		fmt.Printf(" ")
		printBoolExprType(be.Boolop)
	}

	// Drop parentheses when change operators
	parentheses := be.Boolop != prevOp

	if parentheses {
		fmt.Printf("(\n")
		printIndent(1)
	}

	for i := range be.Args.Items {
		if int(be.Boolop) == 2 {
			printBoolExprType(be.Boolop)
		}

		if tbe, ok := be.Args.Items[i].(nodes.BoolExpr); ok {
			printBoolExpr(tbe, be.Boolop)
		} else {
			printNode(be.Args.Items[i])
		}

		if i < len(be.Args.Items)-1 {
			fmt.Println()
			printIndent(1)
			printBoolExprType(be.Boolop)
		}
	}

	if parentheses {
		fmt.Printf("\n)")
	}
}

func printBoolExprType(exprType nodes.BoolExprType) {
	switch exprType {
	case nodes.AND_EXPR:
		fmt.Printf(keywordFormatter("and "))

	case nodes.OR_EXPR:
		fmt.Printf(keywordFormatter("or "))

	default:
		if int(exprType) == 2 {
			fmt.Printf(keywordFormatter("not "))
		}
	}
}

func printAExpr(ae nodes.A_Expr) {
	printNode(ae.Lexpr)

	switch ae.Kind {
	case nodes.AEXPR_OP_ANY:
		printAExprKeywords(ae.Name)
		fmt.Printf("%s(", keywordFormatter("any"))
		printNode(ae.Rexpr)
		fmt.Printf(")")

	case nodes.AEXPR_OP_ALL:
		printAExprKeywords(ae.Name)
		fmt.Printf("%s(", keywordFormatter("all"))
		printNode(ae.Rexpr)
		fmt.Printf(")")

	case nodes.AEXPR_BETWEEN,
		nodes.AEXPR_NOT_BETWEEN:
		printAExprKeywords(ae.Name)
		printNode(ae.Rexpr.(nodes.List).Items[0])
		fmt.Printf(" %s ", keywordFormatter("and"))
		printNode(ae.Rexpr.(nodes.List).Items[1])

	case nodes.AEXPR_LIKE:
		fmt.Printf("%s ", keywordFormatter(" like"))
		printNode(ae.Rexpr)

	case nodes.AEXPR_ILIKE:
		fmt.Printf("%s ", keywordFormatter(" ilike"))
		printNode(ae.Rexpr)

	case nodes.AEXPR_OP:
		// Prints the operator
		printAExprKeywords(ae.Name)

		printNode(ae.Rexpr)

	default:
		fmt.Printf("\nAExpr: %v\n", ae.Kind)
		p()

	}
}

func printAExprKeywords(op nodes.List) {
	for i := range op.Items {
		// These are keywords
		fmt.Printf(" %s ", keywordFormatter(op.Items[i].(nodes.String).Str))
	}
}

func printNullTest(nt nodes.NullTest) {
	if nt.Xpr != nil {
		printNode(nt.Xpr)
	}

	printNode(nt.Arg)

	if nt.Nulltesttype == nodes.IS_NULL {
		fmt.Printf(keywordFormatter(" is null"))
	} else {
		fmt.Printf(keywordFormatter(" is not null"))
	}
}

func printKeyWord(kw string) {
	fmt.Printf(keywordFormatter(kw))
}

func printAConst(ac nodes.A_Const) {
	switch ac.Val.(type) {
	case nodes.String:
		fmt.Printf("'")
		printNode(ac.Val)
		fmt.Printf("'")
	default:
		printNode(ac.Val)
	}
}

func printCommonTableExpr(cte nodes.CommonTableExpr) {
	if cte.Ctename != nil {
		fmt.Printf("%s (\n", *cte.Ctename)
	}

	// TODO indenti-magica
	if cte.Cterecursive {
		p()
	}

	printNode(cte.Ctequery)

	fmt.Printf("\n)\n")
}

func printParamRef(pr nodes.ParamRef) {
	if param, ok := detectedParameters[paramCounter]; ok {
		fmt.Printf("%s", param)
		paramCounter++
	} else {
		fmt.Printf("?")
	}
}

func printNode(node nodes.Node) {
	switch node.(type) {
	case nodes.RawStmt:
		printNode(node.(nodes.RawStmt).Stmt)

	case nodes.SelectStmt:
		printSelectStatement(node.(nodes.SelectStmt))

	case nodes.IntoClause:
		fmt.Printf("into\n")

	case nodes.ResTarget:
		fmt.Printf("%s", printResTarget(node.(nodes.ResTarget)))

	case nodes.ColumnRef:
		fmt.Printf("%s", printColumRef(node.(nodes.ColumnRef)))

	case nodes.A_Star:
		fmt.Printf("*")

	case nodes.String:
		fmt.Printf("%s", node.(nodes.String).Str)

	case nodes.JoinExpr:
		printJoin(false, node.(nodes.JoinExpr))

	case nodes.RangeVar:
		printRangeVar(node.(nodes.RangeVar))

	case nodes.BoolExpr:
		printBoolExpr(node.(nodes.BoolExpr), node.(nodes.BoolExpr).Boolop)

	case nodes.A_Expr:
		printAExpr(node.(nodes.A_Expr))

	case nodes.A_Const:
		printAConst(node.(nodes.A_Const))

	case nodes.Integer:
		fmt.Printf("%d", node.(nodes.Integer).Ival)

	case nodes.NullTest:
		printNullTest(node.(nodes.NullTest))

	case nodes.CommonTableExpr:
		printCommonTableExpr(node.(nodes.CommonTableExpr))

	case nodes.ParamRef:
		printParamRef(node.(nodes.ParamRef))

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

func main() {
	sql := `
	with tab as (
		select shit
		from somewhere
		where something = 42
	), tab2 as (
		select
			*
		from tab5
		cross join some_info
		where tab4.id = ?s
		and tab4.ids = any(?named_parameters_like_a_boss)
		and tab4.others = any('{3,4,5,6}')
	)
	SELECT *, t.id
	from tab t
	join tab2 t2      on t2.id = t.tab_id
	where t.name like 'thing%'
	and t.name = '24'
	and t.num = 33
	and t.arr = any(t2.arr)
	and t.range between 20 and 2000
	and t.range2 not between 20 and 300
	and not t.bval
	and (
		t.opt is not null
		or (
			t.opt is null
			and t.not_opt
		)
	)
`

	workingSQL := processNamedParameters(sql)

	tree, err := pg_query.Parse(workingSQL)

	if err != nil {
		panic(err)
	}

	for i := range tree.Statements {
		printNode(tree.Statements[i])
	}

	fmt.Println()

}
