# pg_flame

A flamegraph generator for Postgres `EXPLAIN ANALYZE` output.

<a href="https://mgartner.github.io/pg_flame/flamegraph.html">
  <img width="700" src="https://user-images.githubusercontent.com/1128750/67529948-7d46b000-f672-11e9-918a-a6e9a1e42cd4.png">
</a>

## Demo

Try the demo [here](https://mgartner.github.io/pg_flame/flamegraph.html).

## Installation

Download one of the compiled binaries [in the releases
tab](https://github.com/mgartner/pg_flame/releases).

If you'd like to build a binary from the source code, run the following
commands. Note that compiling requires Go version 1.13+.

```
git clone https://github.com/mgartner/pg_flame.git
cd pg_flame
go build
```

## Usage

1. Generate a query plan in JSON by prefixing a SQL query with `EXPLAIN
   (ANALYZE, BUFFERS, FORMAT JSON)`. Save the output to a file. Example query
   plan JSON can be found
   [here](https://mgartner.github.io/pg_flame/plan.json).

_Example:_

```
psql lob_local -qAtc 'EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) SELECT id FROM users' > plan.json
```

2. Then generate the flamegraph by passing the JSON as standard input to
`pg_flame` and direct standard output to a file.

_Example:_

```
cat plan.json | ./pg_flame > flamegraph.html
```

3. Open `flamegraph.html` in a browser of your choice.

## Background

[Flamegraphs](http://www.brendangregg.com/flamegraphs.html) were invented by
Brendan Gregg to visualize CPU consumption of profiled code-paths of software.
They are useful visualization tools in many types of performance
investigations. Luca Canali has previously shown the benefits of using
[flamegraph visualizations of Oracle database
profiles](https://externaltable.blogspot.com/2014/05/flame-graphs-for-oracle.html)
for debugging slow database queries.

Pg_flame is in extension of that work for Postgres query plans. Instead of
being used to graph CPU time of internal Postgres functions, it generates a
visual hierarchy of query plans. This visualization identifies the relative
time of each part of a query plan.

This tool relies on the
[`spiermar/d3-flame-graph`](https://github.com/spiermar/d3-flame-graph) plugin to
generate the flamegraph.
