package processors

import (
	"github.com/dbreedt/pgPretty/interfaces"
	pg_query "github.com/pganalyze/pg_query_go"
)

// ProcessSQL Uses the PostgresSQL parser to gain an AST. The AST is then used to start the formatting process
//            by utilising the formatter and printer provided.
func ProcessSQL(sql string, formatter interfaces.PgSqlFormatter) (string, error) {
	tree, err := pg_query.Parse(sql)
	if err != nil {
		return "", err
	}

	for i := range tree.Statements {
		formatter.PrintNode(tree.Statements[i])
	}

	return formatter.String(), nil
}
