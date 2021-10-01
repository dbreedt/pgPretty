select *
from tab7
where t.id in (select id from some_id_store where things = 'borked')
and t.id not in (select id from some_id_store where things = 'borked22')
and t.id = any(select id from some_id_store where things = 'bbb')
and t.id4 = all(select id from some_id_store where things = 'bbb')
and exists (
    select from x222yy where id = t2.id
)
and not exists (
    select from x23423z where id = t2.id
)