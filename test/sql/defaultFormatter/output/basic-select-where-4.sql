{{ .Select}}
{{ .Ws}}*
{{ .From}}
{{ .Ws}}tab7 t7
{{ .Where}}
{{ .Ws}}t7.kids = {{ .Any}}(t7.parents)
{{ .Ws}}{{ .Or}} t7.parents = {{ .All}}(t7.kids)