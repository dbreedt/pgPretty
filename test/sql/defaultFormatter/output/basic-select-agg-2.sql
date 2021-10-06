{{ .Select}}
{{ .Ws}}id,
{{ .Ws}}{{ .Fn "count"}}(*),
{{ .Ws}}{{ .Fn "count"}}(t7.qty)
{{ .Ws}}{{ .Ws}}{{ .Filter}} (
{{ .Ws}}{{ .Ws}}{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}{{ .Ws}}t7.ft = 42
{{ .Ws}}{{ .Ws}}) {{ .As}} "qty",
{{ .Ws}}{{ .Fn "sum"}}(age)
{{ .Ws}}{{ .Ws}}{{ .Filter}} (
{{ .Ws}}{{ .Ws}}{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}{{ .Ws}}t7.tp_type = 22
{{ .Ws}}{{ .Ws}}{{ .Ws}}{{ .Ws}}{{ .And}} t7.tp_type2 = 77
{{ .Ws}}{{ .Ws}}) {{ .As}} "mon_val",
{{ .Ws}}{{ .Fn "max"}}(t7.counter) {{ .As}} "Large",
{{ .Ws}}{{ .Fn "min"}}(t7.counter) {{ .As}} "small",
{{ .Ws}}{{ .Fn "avg"}}(t7.val) {{ .As}} "mean"
{{ .From}}
{{ .Ws}}tab7 t7
{{ .Group}} {{ .By}}
{{ .Ws}}t7.id
{{ .Having}}
{{ .Ws}}3 > 11