package main

import (
	"fmt"

	"github.com/dbreedt/pgPretty/formatters"
	helpers "github.com/dbreedt/pgPretty/helpers"
	printers "github.com/dbreedt/pgPretty/printers"
	pg_query "github.com/pganalyze/pg_query_go"
)

func main() {
	// TODO: Process parameters here

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
	--join lateral (select * from tab8 where tab8_key = t.lame_key) t3 on true
	where not t.id and t.name like 'thing%'
	and t.name = '24'
	and t.num = 33
	and t.numf = 22.321231
	and t.arr = any(t2.arr)
	and t.range between 20 and 2000
	and t.range2 not between 20 and 300
	and not t.bval
	--and t.id in (select id from some_id_store where things = 'borked')
	and (
		t.opt is not null
		or (
			t.opt is null
			and t.not_opt
		)
	)	and x = ?Name
`
	// remove any illegal named parameters and store them for later processing
	workingSQL, detectedParameters := helpers.ProcessNamedParameters(sql)

	tree, err := pg_query.Parse(workingSQL)
	if err != nil {
		panic(err)
	}

	printer := printers.NewSpacePrinter(true, 4)
	formatter := formatters.NewDefaultFormatterWithParameters(printer, detectedParameters)

	for i := range tree.Statements {
		formatter.PrintNode(tree.Statements[i], false)
	}

	fmt.Println(printer)
}
