# pg_flame

A flamegraph generator for Postgres `EXPLAIN ANALYZE` output.

## Usage

1. Generate a query plan in JSON by prefixing a SQL query with
`EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)`. Save the output to a file.

_Example:_

```
psql lob_local -qAtc 'EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) SELECT id FROM users' > plan.json
```

2. Then generate the flamegraph by passing the JSON as standard input to
`pg_flame` and direct standard output to a file.

_Example:_

```
cat plans.json | ./pg_flame > flamegraph.html
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
