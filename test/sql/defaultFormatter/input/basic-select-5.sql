select *
from generate_series(1,100) as g(d)
join tab2 t2 on t2.day = g.d