{{ .Select}}
{{ .Ws}}*
{{ .From}}
{{ .Ws}}tab7 t
{{ .Where}}
{{ .Ws}}{{ .Not}} t.id
{{ .Ws}}{{ .And}} t.name {{ .Like}} 'thing%'
{{ .Ws}}{{ .And}} t.name = '24'
{{ .Ws}}{{ .And}} t.num = 33
{{ .Ws}}{{ .And}} t.numf = 22.321231
{{ .Ws}}{{ .And}} t.arr = {{ .Any}}(t2.arr)
{{ .Ws}}{{ .And}} t.range {{ .Between}} 20 {{ .And}} 2000
{{ .Ws}}{{ .And}} t.range2 {{ .Not}} {{ .Between}} 20 {{ .And}} 300
{{ .Ws}}{{ .And}} {{ .Not}} t.bval
{{ .Ws}}{{ .And}} (
{{ .Ws}}{{ .Ws}}t.opt {{ .Is}} {{ .Not}} {{ .Null}}
{{ .Ws}}{{ .Ws}}{{ .Or}} (
{{ .Ws}}{{ .Ws}}{{ .Ws}}t.opt {{ .Is}} {{ .Null}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}{{ .And}} t.not_opt
{{ .Ws}}{{ .Ws}})
{{ .Ws}})
{{ .Ws}}{{ .And}} x = 22
{{ .Ws}}{{ .And}} t.name {{ .Ilike}} '%bob%'