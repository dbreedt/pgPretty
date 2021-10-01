{{ .Select}}
{{ .Ws}}*
{{ .From}}
{{ .Ws}}tab7
{{ .Where}}
{{ .Ws}}t.id {{ .In}}(
{{ .Ws}}{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}id
{{ .Ws}}{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}some_id_store
{{ .Ws}}{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}things = 'borked'
{{ .Ws}})
{{ .Ws}}{{ .And}} {{ .Not}} t.id {{ .In}}(
{{ .Ws}}{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}id
{{ .Ws}}{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}some_id_store
{{ .Ws}}{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}things = 'borked22'
{{ .Ws}})
{{ .Ws}}{{ .And}} t.id = {{ .Any}}(
{{ .Ws}}{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}id
{{ .Ws}}{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}some_id_store
{{ .Ws}}{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}things = 'bbb'
{{ .Ws}})
{{ .Ws}}{{ .And}} t.id4 = {{ .All}}(
{{ .Ws}}{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}id
{{ .Ws}}{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}some_id_store
{{ .Ws}}{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}things = 'bbb'
{{ .Ws}})
{{ .Ws}}{{ .And}} {{ .Exists}}(
{{ .Ws}}{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}x222yy
{{ .Ws}}{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}id = t2.id
{{ .Ws}})
{{ .Ws}}{{ .And}} {{ .Not}} {{ .Exists}}(
{{ .Ws}}{{ .Ws}}{{ .Select}}
{{ .Ws}}{{ .Ws}}{{ .From}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}x23423z
{{ .Ws}}{{ .Ws}}{{ .Where}}
{{ .Ws}}{{ .Ws}}{{ .Ws}}id = t2.id
{{ .Ws}})