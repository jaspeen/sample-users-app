# Simple Golang graphql backend with postgres persistence

See [api/schema.graphqls](api/schema.graphqls) for GraphQL definition

Will automatically created db, tables and admin user on start

## Build

Requred Go 1.18+

```sh
go build cmd/server/server.go
```

## Run

```sh
export SMPL_DBCONNECT="sampleuser:samplepassword@localhost/sampledb?sslmode=disable"; ./server
```

See `server --help` for list of available env vars