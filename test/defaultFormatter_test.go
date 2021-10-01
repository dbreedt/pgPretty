package test

import (
	"bytes"
	"io/ioutil"
	"path"
	"strings"
	"testing"
	"text/template"

	"github.com/dbreedt/pgPretty/formatters"
	"github.com/dbreedt/pgPretty/printers"
	"github.com/dbreedt/pgPretty/processors"
	"github.com/kylelemons/godebug/pretty"
)

const (
	baseDir = "./sql/defaultFormatter/"
	filter  = "" // for debugging purposes: add the name of the sql file you want to test to avoid all the others from being tested
)

func TestSqlFiles(t *testing.T) {
	files, err := ioutil.ReadDir(path.Join(baseDir, "input"))
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if filter != "" && file.Name() != filter {
			continue
		}

		spacePrinter := printers.NewBasePrinter(false, true, 2)
		dfSpace := formatters.NewDefaultFormatter(spacePrinter)

		srcData, err := ioutil.ReadFile(path.Join(baseDir, "input", file.Name()))
		if err != nil {
			t.Log("intput", file.Name(), err)
			t.Fail()
		}

		sqlOut, err := processors.ProcessSQL(string(srcData), dfSpace)
		if err != nil {
			t.Log(file.Name(), err)
			t.Fail()
		}

		expectedData, err := ioutil.ReadFile(path.Join(baseDir, "output", file.Name()))
		if err != nil {
			t.Log("output", file.Name(), err)
			t.Fail()
		}

		sqlExp, err := processKeywords(string(expectedData), true, false)
		if err != nil {
			t.Fatal("keywords", file.Name(), err)
		}

		if sqlOut != sqlExp {
			t.Log("SPACE", file.Name())
			t.Log(pretty.Compare(strings.Split(sqlOut, "\n"), strings.Split(sqlExp, "\n")))
			t.Fail()
		}

		tabPrinter := printers.NewBasePrinter(true, false, 1)
		dfTab := formatters.NewDefaultFormatter(tabPrinter)

		sqlOut, err = processors.ProcessSQL(string(srcData), dfTab)
		if err != nil {
			t.Log(file.Name(), err)
			t.Fail()
		}

		sqlExp, err = processKeywords(string(expectedData), false, true)
		if err != nil {
			t.Fatal("keywords", file.Name(), err)
		}

		if sqlOut != sqlExp {
			t.Log("TAB", file.Name())
			t.Log(pretty.Compare(strings.Split(sqlOut, "\n"), strings.Split(sqlExp, "\n")))
			t.Fail()
		}
	}
}

func processKeywords(sql string, upper bool, tabs bool) (string, error) {
	templ, err := template.New("keywords").Parse(sql)
	if err != nil {
		return "", err
	}

	output := bytes.Buffer{}
	kw := NewKeywords(upper)
	if tabs {
		kw.Ws = "\t"
	} else {
		kw.Ws = "  "
	}

	err = templ.Execute(&output, kw)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
