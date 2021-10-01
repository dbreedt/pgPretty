{{ .Select}}
{{ .Ws}}{{ .Distinct}} a,
{{ .Ws}}b,
{{ .Ws}}c,
{{ .Ws}}d
{{ .From}}
{{ .Ws}}tab7 t7
{{ .Join}}
{{ .Ws}}tab8 t8
{{ .Ws}}{{ .On}}
{{ .Ws}}{{ .Ws}}t8.id = t7.id
{{ .Left}} {{ .Join}}
{{ .Ws}}tab9 t9
{{ .Ws}}{{ .On}}
{{ .Ws}}{{ .Ws}}t9.id = t8.id