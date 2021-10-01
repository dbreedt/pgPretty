{{ .With}} data {{ .As}} (
{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}*
{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}tab22
{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}something > 42
{{ .Ws}}{{ .Limit}}
{{ .Ws}}{{ .Ws}}1
)
{{ .Select}}
{{ .Ws}}t7.*,
{{ .Ws}}t8.id,
{{ .Ws}}t8.k1,
{{ .Ws}}t9.*
{{ .From}}
{{ .Ws}}tab7 t7
{{ .Cross}} {{ .Join}}
{{ .Ws}}data
{{ .Join}}
{{ .Ws}}tab8 t8
{{ .Ws}}{{ .On}}
{{ .Ws}}{{ .Ws}}t8.id = t7.id
{{ .Left}} {{ .Join}}
{{ .Ws}}tab9 t9
{{ .Ws}}{{ .On}}
{{ .Ws}}{{ .Ws}}t9.id = t8.id
{{ .Left}} {{ .Join}}
{{ .Ws}}(
{{ .Ws}}{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}*
{{ .Ws}}{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}tab1
{{ .Ws}}{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}x = 'y'
{{ .Ws}}) t1
{{ .Ws}}{{ .On}}
{{ .Ws}}{{ .Ws}}t1.id = t7.id
{{ .Left}} {{ .Join}} {{ .Lateral}}
{{ .Ws}}(
{{ .Ws}}{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}*
{{ .Ws}}{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}tab2
{{ .Ws}}{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}id = t7.id
{{ .Ws}}) t2
{{ .Ws}}{{ .On}}
{{ .Ws}}{{ .Ws}}true