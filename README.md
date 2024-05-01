# d8b

d8b is a TUI built to query and manipulate postgres databases. Built using the [tview](https://github.com/rivo/tview/tree/master) framework.

The aim is to store popular/saved queries and provide these as fast options. As well as running new queries and displaying that data. I currently use `psql` alot and I am building this to make that process better/faster/happier!

## Configuration

- create a `config.toml` file in the root directory, as follows:

```toml
host = "localhost"
port = 5432
user = "postgres"
password = "someSecureAsPassword"
dbname = "postgres"
```

## TODO

- [ ] Deal with tenanted tables
- [ ] Edit tables
- [ ] Multiple configs (for multiple databases)
- [ ] Proper errors for bad queries
- [ ] Create executable
- [ ] Add tests? + other go project files
