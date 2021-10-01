# pgPretty
Go tool to beautify PostgreSQL statements

Needs a bit of work still


### Build
```bash
make build
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
The test uses input files from `./sql/input` runs the content through a formatter and then compares the output from the formatter with that of the file with the same name in `./sql/output`.

The output files use `text/template` syntax to reduce the amount of output files that need to be created when testing different formatting rules (spaces, tabs, upper case keywords, etc.)

### Todo
* convert to cmd structure
* add support for missing sql syntax
* other formatters
