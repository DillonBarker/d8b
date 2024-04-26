# d8b

d8b is a TUI built to query and manipulate postgres databases. Built using the bubbletea framework.

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

## Queries

- Use the add **new query** functionality or, create a `queries.toml` file in the root directory, as follows:

```toml
[[choice]]
name = "List users"
query = "SELECT * FROM user"

[[choice]]
name = "List user permissions"
query = "SELECT * FROM dillon.permissions"
```

## TODO

- [ ] Simplify existing code and split into packages, maybe seperate update, view, render if possible.
- [ ] Create new queries, edit queries and delete queries
- [ ] Enter on table gives extra information
- [ ] Multiple configs (for multiple databases)
- [ ] Proper errors for bad queries
- [ ] Create executable
- [ ] Add tests? + other go project files
