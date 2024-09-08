## How to start Dev

1. Setup Postgres:16 database

```bash
$ docker-compose up -d
```

2. Run goose up

```bash
$ cd migrations
$ goose postgres "postgres://postgres:postgres@localhost:5432/iws" up
$ cd ..
```
