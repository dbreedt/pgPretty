with data as (
    select * from tab
), other_data as (
    select * from tab2
)
select *
from data d
join other_data od on od.id = d.id