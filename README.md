# d8b

d8b is a TUI built to query and manipulate postgres databases. Built using the [tview](https://github.com/rivo/tview/tree/master) framework.

The aim is to display database data in a similar way to my favourite kubernetes TUI k9s. I currently use `psql` alot and I am building this to make that process better/faster/happier, and to practice go.

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
