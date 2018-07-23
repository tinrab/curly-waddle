# GraphQL Server Using PostgreSQL Window Functions

Code for the [N+1 Queries in GraphQL Using PostgreSQL Window Functions](https://outcrawl.com/graphql-postgresql-window-functions) article.

Start the database:

```
$ docker-compose up -d
```

Optionally insert fake data:

```
$ vgo run ./cmd/fakedata/main.go
```

Start GraphQL server:

```
$ vgo run .
```

See [schema.graphql](./graphql/schema.graphql) for possible queries.
