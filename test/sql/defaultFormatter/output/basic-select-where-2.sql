{{ .Select}}
{{ .Ws}}*
{{ .From}}
{{ .Ws}}tab7 t
{{ .Where}}
{{ .Ws}}!t.id
{{ .Ws}}{{ .And}} t.name *~ 'thing'
{{ .Ws}}{{ .And}} (
{{ .Ws}}{{ .Ws}}t.opt {{ .Is}} {{ .Null}}
{{ .Ws}}{{ .Ws}}{{ .Or}} (
{{ .Ws}}{{ .Ws}}{{ .Ws}}t.opt {{ .Is}} {{ .Not}} {{ .Null}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}{{ .And}} t.no_opt
{{ .Ws}}{{ .Ws}})
{{ .Ws}})
{{ .Ws}}{{ .And}} (
{{ .Ws}}{{ .Ws}}(
{{ .Ws}}{{ .Ws}}{{ .Ws}}t.opt2
{{ .Ws}}{{ .Ws}}{{ .Ws}}{{ .And}} t.opt7
{{ .Ws}}{{ .Ws}})
{{ .Ws}}{{ .Ws}}{{ .Or}} (
{{ .Ws}}{{ .Ws}}{{ .Ws}}t.opt3
{{ .Ws}}{{ .Ws}}{{ .Ws}}{{ .And}} t.opt4
{{ .Ws}}{{ .Ws}})
{{ .Ws}})