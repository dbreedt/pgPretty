{{ .Select}}
{{ .Ws}}t7.*,
{{ .Ws}}(
{{ .Ws}}{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}name
{{ .Ws}}{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}people p
{{ .Ws}}{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}p.pers_no = t7.person_id
{{ .Ws}}) fname,
{{ .Ws}}t7.x + t7.y score,
{{ .Ws}}{{ .Exists}}(
{{ .Ws}}{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}fired
{{ .Ws}}{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}pers_no = t7.person_id
{{ .Ws}}) fired
{{ .From}}
{{ .Ws}}tab7 t7