package main

import (
	"fmt"
	"strings"

	pg_query "github.com/lfittl/pg_query_go"
)

type caseFormatter func(s string) string

var indent = "  "
var keywordFormatter caseFormatter = strings.ToLower
var detectedParameters = make(map[int]string)
var paramCounter int
var formatter *Formatter

func p() {
	panic("not supported")
}

func main() {

	// TODO: Process parameters here

	formatter = NewFormatter(false, true, 2)

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
	where not t.id and t.name like 'thing%'
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
		printNode(tree.Statements[i], false)
	}

	fmt.Println()

}
