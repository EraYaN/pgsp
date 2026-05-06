# pgsp - PostgreSQL Stat Progress CLI Monitor

> **⚠️ This repository is a fork of [noborus/pgsp](https://github.com/noborus/pgsp).**

[![Go Reference](https://pkg.go.dev/badge/github.com/EraYaN/pgsp.svg)](https://pkg.go.dev/github.com/EraYaN/pgsp)

A TUI tool that monitors PostgreSQL's pg_stat_progress*. This also resolves `oid` columns by connecting to the relevant databases and accessing `pg_catalog.pg_class`.

Supported progress reports are ANALYZE, CLUSTER, CREATE INDEX, VACUUM, COPY, and BASE_BACKUP.
See [Progress Reporting](https://www.postgresql.org/docs/current/progress-reporting.html) for more information.

![pgsp.png](https://raw.githubusercontent.com/EraYaN/pgsp/master/docs/pgsp.png)

## Requires

go 1.25 or later

## Install

### Download binary

[releases page](https://github.com/EraYaN/pgsp/releases/).

### Go install

```console
go install github.com/EraYaN/pgsp/cmd/pgsp@latest
```

## Usage

Shows a progress bar if pg_stat_progress* is updated while waiting while running. 

```console
$ pgsp --dsn 'host=/var/run/postgresql port=5432'
 Monitor: Analyze BaseBackup Cluster Copy CreateIndex Vacuum, Connections: 2
 quit: q, ctrl+c, esc; scroll: arrow keys

pg_stat_progress_cluster: benchmark, pgbench_accounts
Command: VACUUM FULL; Phase: seq scanning heap; PID: 2605557; RELID: 997696163;
Heap Tuples Scanned: 33611488; Written: 33611488
Heap Blocks: 551009/1639345 (33.6%)
▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  34%
```

It is also possible to specify one of the `analyze`, `basebackup`, `cluster`, `createindex`, `vacuums`, `copy` for monitoring.

```console
pgsp basebackup
```

```console
Monitors PostgreSQL's pg_stat_progress_*.
Analyze, BaseBackup, Cluster, Copy, CreateIndex, Vacuum can be specified.

Usage:
  pgsp [flags]

Flags:
  -a, --AfterCompletion int   Time to display after completion(Seconds) (default 10)
  -i, --Interval float        Update interval(Seconds) (default 0.5)
      --config string         config file (default is $HOME/.pgsp.yaml)
      --debug                 debug message for toggle
      --dsn string            PostgreSQL data source name
  -f, --fullscreen            Display in Full Screen
  -h, --help                  help for pgsp
  -t, --toggle                Help message for toggle
  -v, --version               display version information
```
