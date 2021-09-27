package interfaces

import nodes "github.com/pganalyze/pg_query_go/nodes"

/*
PgSqlFormatter Basic interface for a postgres AST formatter
*/
type PgSqlFormatter interface {
	PrintNode(node nodes.Node)
	String() string
}
