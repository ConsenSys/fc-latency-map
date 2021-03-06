# FC Latency Map - Manager

## Description

FC Latency Map - Manager allows to get latency measurements

## Requirements

### Packages

- swagger

```shell
go get -u github.com/go-swagger/go-swagger/cmd/swagger

```

- sqlite3

```bash
sudo apt install sqlite3
```

### Configuration

Must have an ".env" file with the configurations to run.

Is available an example of configuration on [.env.example](./.env.example)

```bash
cp .env.example .env
```

| Key                           | Value type | Description                                                                                                            |
| ----------------------------- | ---------- | ---------------------------------------------------------------------------------------------------------------------- |
| SERVICE_NAME                  | string     | Service name.                                                                                                          |
| FILECOIN_NODE_URL             | string     | Lotus Filecoin node url.                                                                                               |
| FILECOIN_BLOCKS_OFFSET        | string     | Filecoin blocks search offset.                                                                                         |
| DB_CONNECTION                 | string     | SQLite file path.                                                                                                      |
| RIPE_API_KEY                  | string     | [Ripe Atlas API Key](https://atlas.ripe.net/keys/)                                                                     |
| RIPE_ONE_OFF                  | boolean    | On ping only to get latency. When is 'true' the RIPE_TRACEROUTE_INTERVAL and RIPE_TRACEROUTE_RUNNING_TIME are ignored. |
| RIPE_TRACEROUTE_INTERVAL      | number     | Interval between ping to get miners latency.                                                                           |
| RIPE_TRACEROUTE_RUNNING_TIME  | number     | Running period get latency.                                                                                            |
| RIPE_PACKETS                  | number     | Amount of packets to ping miners.                                                                                      |
| RIPE_REQUESTED_PROBES         | number     | Max number of probes to call miners.                                                                                   |
| RIPE_PROBES_PER_AIRPORT       | number     | Number of probes closest to airports.                                                                                  |
| NEAREST_AIRPORTS              | number     | Number of airports closest to miner.                                                                                   |
| CONSTANT_AIRPORTS             | string     | Airports JSON file path.                                                                                               |
| SQL_DEBUG                     | boolean    | Display SQL debug.                                                                                                     |
| CRON_SCHEDULE_CREATE_MEASURES | string     | Cron time to export measures.                                                                                          |
| CRON_SCHEDULE_IMPORT_MEASURES | string     | Cron time to import measures.                                                                                          |
| WEBHOOK_NOTIFY_URLS           | string     | Webhooks links.                                                                                                        |

## Development

- Run golangci

```shell
# run golint-ci
golangci-lint run ./... --fix
```

- Install pre-commit hooks

```shell
# @ project root
pre-commit install
```

- Execute pre-commit manually

```shell
# @ project root
pre-commit run --all-files
```

## SQLite commands

1. Open database

```bash
sqlite3 data/database.db

SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
sqlite>
```

2. View tables

```bash
sqlite> .tables

>>> miners
```

3. Select miners

```bash
sqlite> select * from miners;

>>> 1|2021-09-07 17:04:53.41480159+02:00|2021-09-07 17:04:53.41480159+02:00||dummyAddress|dummyIp
```

4. Close database

```bash
sqlite> .quit
```

or

```bash
sqlite> ^C^C^C
```

## Command-line Interface

### Usage

1. Run CLI

```bash
go run cmd/cli/main.go
```

2. Update miners list
   without offset (it will apply by default the offset from the .env file)

```bash
>>> miners-update
```

or with an offset of 10 (latest block heights)

```bash
>>> miners-update 10
```

3. Parse miners from a given block height

```bash
>>> miners-parse 1107781
```

4. List miners

```bash
>>> miners-list
```

5. Load all the miner from MarketDeals

```bash
>>> miners-parse-state-market
```

### Demo

```shell
rm data/database.db
go run cmd/cli/main.go miners-parse-state-market
go run cmd/cli/main.go locations-update large
go run cmd/cli/main.go probes-import
go run cmd/cli/main.go probes-update

go run cmd/cli/main.go measures-create
#   wait until have ripe results
go run cmd/cli/main.go measures-get
go run cmd/cli/main.go measures-export
```

## API

The Manager exposes an API to allow health check and metrics request.

### Health Check

Open [http://localhost:3001/health-check](http://localhost:3001/health-check)

It should respond:

```
{
"success": true
}
```

### Metrics

Open [http://localhost:3001/metrics](http://localhost:3001/metrics)

It should respond:

```
{
"locations": "606",
"miners": "140",
"probes": "1890"
}
```

### Swagger

[./swagger.yml](./swagger.yml)
