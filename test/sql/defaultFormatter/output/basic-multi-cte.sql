{{ .With}} data {{ .As}} (
{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}*
{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}tab
),
other_data {{ .As}} (
{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}*
{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}tab2
)
{{ .Select}}
{{ .Ws}}*
{{ .From}}
{{ .Ws}}data d
{{ .Join}}
{{ .Ws}}other_data od
{{ .Ws}}{{ .On}}
{{ .Ws}}{{ .Ws}}od.id = d.id