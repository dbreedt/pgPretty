select *
from tab7 t
where !t.id
and t.name *~ 'thing'
and (
    t.opt is null
    or (t.opt is not null and t.no_opt)
)
and (
    (t.opt2 and t.opt7)
    or (t.opt3 and t.opt4)
)