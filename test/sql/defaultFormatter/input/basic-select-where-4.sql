select *
from tab7 t7
where t7.kids = any(t7.parents)
or t7.parents = all(t7.kids)