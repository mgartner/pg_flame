# pg_flame [![Version](https://img.shields.io/badge/version-v1.1-blue.svg)](https://github.com/mgartner/pg_flame/releases) [![Build Status](https://travis-ci.com/mgartner/pg_flame.svg?branch=master)](https://travis-ci.com/mgartner/pg_flame)

A flamegraph generator for Postgres `EXPLAIN ANALYZE` output.

<a href="https://mgartner.github.io/pg_flame/flamegraph.html">
  <img width="700" src="https://user-images.githubusercontent.com/1128750/67738754-16f0c300-f9cd-11e9-8fc2-6acc6f288841.png">
</a>

## Demo

Try the demo [here](https://mgartner.github.io/pg_flame/flamegraph.html).

## Installation

### Download pre-compiled binary

Download one of the compiled binaries [in the releases
tab](https://github.com/mgartner/pg_flame/releases). Once downloaded, move
`pg_flame` into your `$PATH`.

### Docker

Alternatively, if you'd like to use Docker to build the program, you can.

```
$ git clone https://github.com/mgartner/pg_flame.git
$ cd pg_flame
$ docker build --tag 'pg_flame' .
```

### Build from source

If you'd like to build a binary from the source code, run the following
commands. Note that compiling requires Go version 1.13+.

```
$ git clone https://github.com/mgartner/pg_flame.git
$ cd pg_flame
$ go build
```

A `pg_flame` binary will be created that you can place in your `$PATH`.

## Usage

The `pg_flame` program reads a JSON query plan from standard input and writes
the flamegraph HTML to standard ouput. Therefore you can pipe and direct input
and output however you desire.

### Example: One-step

```bash
$ psql dbname -qAtc 'EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) SELECT id FROM users' \
    | pg_flame \
    > flamegraph.html \
    && open flamegraph.html
```

### Example: Multi-step with SQL file

Create a SQL file with the `EXPLAIN ANALYZE` query.

```sql
-- query.sql
EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)
SELECT id
FROM users
```

Then run the query and save the JSON to a file.

```bash
$ psql dbname -qAtf query.sql > plan.json
```

Finally, generate the flamegraph HTML.

```
$ cat plan.json | pg_flame > flamegraph.html
```

### Example: Docker

If you've followed the Docker installation steps above, you can pipe query plan JSON to a container and save the output HTML.

```
$ psql dbname -qAtc 'EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) SELECT id FROM users' \
    | docker run -i pg_flame \
    > flamegraph.html
```

## Background

[Flamegraphs](http://www.brendangregg.com/flamegraphs.html) were invented by
Brendan Gregg to visualize CPU consumption per code-path of profiled software.
They are useful visualization tools in many types of performance
investigations. Flamegraphs have been used to visualize Oracle database
[query
plans](https://blog.tanelpoder.com/posts/visualizing-sql-plan-execution-time-with-flamegraphs/)
and [query
executions](https://externaltable.blogspot.com/2014/05/flame-graphs-for-oracle.html)
, proving useful for debugging slow database queries.

Pg_flame is in extension of that work for Postgres query plans. It generates a
visual hierarchy of query plans. This visualization identifies the relative
time of each part of a query plan.

This tool relies on the
[`spiermar/d3-flame-graph`](https://github.com/spiermar/d3-flame-graph) plugin to
generate the flamegraph.
