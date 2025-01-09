# db-stress

Tool to run db stress tests.

### Quickstart

Rename `config-example.json` to `config.json`. Edit the config file as required, add tests, etc.

Build

```bash
go build
```

Run

```bash
./db-stress
```

By default, loads the config from `config.json`. Can also pass config file name as an argument.

```bash
./db-stress -c anotherconfig.json
```

### Config

```json
{
  "label": "Testing DB1", // label for the tests, optional
  "mode": "series", // test running mode, optional, default: series
  "log": false, // whether to log the results to a csv file, optional, default: false
  "connection": {
    "provider": "mssql", // database provider, required
    "connectionString": "sqlserver://username:password@localhost?database=dbname" // database connection string, required
  },
  "tests": [
    {
      "query": "SELECT * FROM users", // sql query, required
      "iterations": 100, // number of times query is run, required
      "workers": 1 // number of go routines running the query concurrently, optional, default: 1
    }
  ]
}
```

Currently supported providers are `mssql` and `postgres`. But adding support for another database would be as simple as adding a driver to the imports.

If the `mode` config is in `series`, the tests will run one after the other. If it is in `parallel`, the tests will run concurrently.
