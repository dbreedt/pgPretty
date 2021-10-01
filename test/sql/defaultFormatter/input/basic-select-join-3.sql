with data as (
    select *
    from tab22
    where something > 42
    limit 1
)
select t7.*, t8.id, t8.k1, t9.*
from tab7 t7
cross join data
join tab8 t8
on t8.id = t7.id
left join tab9 t9
on t9.id = t8.id
left join (select * from tab1 where x = 'y') t1 on t1.id = t7.id
left join lateral (select * from tab2 where id = t7.id ) t2 on true
