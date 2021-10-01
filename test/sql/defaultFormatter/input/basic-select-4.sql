select t7.*,
(select name from people p where p.pers_no = t7.person_id) fname,
t7.x + t7.y score,
exists( select from fired where pers_no = t7.person_id) fired
from tab7 t7