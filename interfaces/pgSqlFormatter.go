package interfaces

import nodes "github.com/pganalyze/pg_query_go/nodes"

/*
PgSqlFormatter Basic interface for a postgres AST formatter
*/
type PgSqlFormatter interface {
	PrintWithClause(wc nodes.WithClause)
	PrintSelectStatement(ss nodes.SelectStmt)
	PrintResTarget(nt nodes.ResTarget, withIndent bool)
	PrintColumnRef(cr nodes.ColumnRef, withIndent bool)
	PrintJoin(first bool, join nodes.JoinExpr)
	PrintAlias(alias *nodes.Alias)
	PrintJoinType(joinType nodes.JoinType, withIndent bool)
	PrintRangeVar(rv nodes.RangeVar, withIndent bool)
	PrintBoolExpr(be nodes.BoolExpr, prevOp nodes.BoolExprType, withIndent bool)
	PrintBoolExprType(exprType nodes.BoolExprType, withIndent bool)
	PrintAExpr(ae nodes.A_Expr, withIndent bool)
	PrintAExprKeywords(op nodes.List)
	PrintNullTest(nt nodes.NullTest, withIndent bool)
	PrintAConst(ac nodes.A_Const)
	PrintCommonTableExpr(cte nodes.CommonTableExpr)
	PrintParamRef(pr nodes.ParamRef, withIndent bool)
	PrintNode(node nodes.Node, withIndent bool)
}
