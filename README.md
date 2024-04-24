# d8b

d8b is a TUI built to query and manipulate postgres databases.

The aim is to store popular/saved queries and provide these as fast options. As well as running new queries and displaying that data.

## Configuration

- create a `config.toml` file in the root directory, as follows:

```toml
host = "localhost"
port = 5432
user = "postgres"
password = "someSecureAsPassword"
dbname = "postgres"
```
