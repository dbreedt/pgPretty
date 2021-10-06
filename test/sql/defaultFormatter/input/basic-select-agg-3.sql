select t7.g1, t7.g2, count(1)
from tab7 t7
group by t7.g1, t7.g2
order by t7.g2 desc, t7.g1 nulls last