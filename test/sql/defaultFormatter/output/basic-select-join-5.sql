{{ .Select}}
{{ .Ws}}t7.*,
{{ .Ws}}t4.id,
{{ .Ws}}t5.id,
{{ .Ws}}t6.id,
{{ .Ws}}t8.id
{{ .From}}
{{ .Ws}}tab7 t7
{{ .Right}} {{ .Join}}
{{ .Ws}}tab4 t4
{{ .Ws}}{{ .On}}
{{ .Ws}}{{ .Ws}}t4.id = t7.id
{{ .Left}} {{ .Join}}
{{ .Ws}}tab5 t5
{{ .Ws}}{{ .On}}
{{ .Ws}}{{ .Ws}}t5.id = t7.id
{{ .Right}} {{ .Join}}
{{ .Ws}}tab6 t6
{{ .Ws}}{{ .On}}
{{ .Ws}}{{ .Ws}}t6.id = t7.id
{{ .Full}} {{ .Join}}
{{ .Ws}}tab8 t8
{{ .Ws}}{{ .On}}
{{ .Ws}}{{ .Ws}}t8.id = t7.id