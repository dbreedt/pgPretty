{{ .Select}}
{{ .Ws}}*
{{ .From}}
{{ .Ws}}{{ .Fn "generate_series"}}(1, 100) g(d)
{{ .Join}}
{{ .Ws}}tab2 t2
{{ .Ws}}{{ .On}}
{{ .Ws}}{{ .Ws}}t2.day = g.d