select
    id,
    count(*),
    count(t7.qty) filter (where t7.ft = 42) qty,
    sum(age) filter (where t7.tp_type = 22 and t7.tp_type2 = 77) mon_val,
    max(t7.counter) as "Large",
    min(t7.counter) "small",
    avg(t7.val) mean
from tab7 t7
group by t7.id
having 3 > 11