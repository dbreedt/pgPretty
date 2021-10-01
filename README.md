# pgPretty
Go tool to beautify PostgreSQL statements

## What does it actually do?
It converts valid, but ugly sql
```sql
--before
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
```

into something a tad more readable

```sql
--after
WITH data AS (
  SELECT
    *
  FROM
    tab22
  WHERE
    something > 42
  LIMIT
    1
)
SELECT
  t7.*,
  t8.id,
  t8.k1,
  t9.*
FROM
  tab7 t7
CROSS JOIN
  data
JOIN
  tab8 t8
  ON
    t8.id = t7.id
LEFT JOIN
  tab9 t9
  ON
    t9.id = t8.id
LEFT JOIN
  (
    SELECT
      *
    FROM
      tab1
    WHERE
      x = 'y'
  ) t1
  ON
    t1.id = t7.id
LEFT JOIN LATERAL
  (
    SELECT
      *
    FROM
      tab2
    WHERE
      id = t7.id
  ) t2
  ON
    true
```

### Usage
```bash
./pgPretty --help
Usage of ./pgPretty:
  -f string
        name of the sql file you want formatted
  -i int
        how many tabs/spaces to use for a single indent (default 2) (default 2)
  -t    use tabs instead of spaces (default is spaces)
  -u    use upper case keywords (default is lower case)
```

### Build
```bash
make build
```

### Test
```bash
make test
```

### Coverage
This assumes you have google-chrome
```bash
make cover
```

### Tests
The test uses input files from `./sql/defaultFormatter/input` runs the content through a formatter and then compares the output from the formatter with that of the file with the same name in `./sql/defaultFormatter/output`.

The output files use `text/template` syntax to reduce the amount of output files that need to be created when testing different formatting rules (spaces, tabs, upper case keywords, etc.)

## Todo
Some items of work still remain
* convert to cmd structure
* add support for missing sql syntax
  * aggregates
  * window functions
  * etc.
* other formatters
