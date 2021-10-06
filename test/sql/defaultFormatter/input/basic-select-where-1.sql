select *
from tab7 t
where not t.id and t.name like 'thing%'
and t.name = '24'
and t.num = 33
and t.numf = 22.321231
and t.arr = any(t2.arr)
and t.range between 20 and 2000
and t.range2 not between 20 and 300
and not t.bval
and (
    t.opt is not null
    or (t.opt is null and t.not_opt)
)
and x = 22
and t.name ilike '%bob%'