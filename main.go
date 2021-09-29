package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dbreedt/pgPretty/formatters"
	helpers "github.com/dbreedt/pgPretty/helpers"
	printers "github.com/dbreedt/pgPretty/printers"
	"github.com/dbreedt/pgPretty/processors"
)

func main() {
	var (
		fileName        string
		useTabs         bool
		capsKeywords    bool
		numIndentations int
	)
	flag.StringVar(&fileName, "file", "", "name of the sql file you want formatted")
	flag.BoolVar(&useTabs, "tabs", false, "use tabs instead of spaces (default is spaces)")
	flag.BoolVar(&capsKeywords, "capsKw", false, "use upper case keywords (default is lower case)")
	flag.IntVar(&numIndentations, "indents", 2, "how many tabs/spaces to use for a single indent (default 2)")

	flag.Parse()

	sql := `
	with tab as (
		select stuff
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
	join tab2 t2      on t2.id = t.tab_id and t2.id != 33
	left join tab7 t7 on t7.id  = t2.id
	right join tab44 t44 on t44.id = t2.id
	left outer join tab8 t8 on t8.id = t2.id
	right outer join tab444 t444 on t444.id = t2.id
	full outer join tab777 t777 on t777.id = t2.id
	left join (select id from x22w ) t99 on t99.id = t2.id
	join lateral (select * from tab8 where tab8_key = t.lame_key) t3 on true
	where not t.id and t.name like 'thing%'
	and t.name = '24'
	and t.num = 33
	and t.numf = 22.321231
	and t.arr = any(t2.arr)
	and t.range between 20 and 2000
	and t.range2 not between 20 and 300
	and not t.bval
	and t.id in (select id from some_id_store where things = 'borked')
	and (
		t.opt is not null
		or (
			t.opt is null
			and t.not_opt
		)
	)	and x = ?Name
	and exists (
		select from x222yy where id = t2.id
	)
	and not exists (
		select from x23423z where id = t2.id
	)
`

	if fileName != "" {
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Println("Failed to open file", fileName, err)
			os.Exit(1)
		}

		sql = string(data)
	}

	// remove any illegal named parameters and store them for later processing
	workingSQL, detectedParameters := helpers.ProcessNamedParameters(sql)
	printer := printers.NewBasePrinter(useTabs, capsKeywords, numIndentations)
	formatter := formatters.NewDefaultFormatterWithParameters(printer, detectedParameters)

	prettySql, err := processors.ProcessSQL(workingSQL, formatter)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(prettySql)
}
