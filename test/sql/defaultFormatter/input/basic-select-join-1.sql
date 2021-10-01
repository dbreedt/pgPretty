select distinct a,b,c,d
from tab7 t7
join tab8 t8
on t8.id = t7.id
left join tab9 t9
on t9.id = t8.id