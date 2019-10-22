# pg_flame

A flamegraph generator for Postgres `EXPLAIN ANALYZE` output.

## Usage

Generate an annotated query plan in JSON.

```
psql lob_local -qAtc 'EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) SELECT id FROM bears' > plan.json
```

Generate the flamegraph.

```
cat plans.json | ./pg_flame > flamegraph.html
```

Open `flamegraph.html` in a browser of your choice.

## Background

[Flamegraphs](http://www.brendangregg.com/flamegraphs.html) were invented by
Brendan Gregg to visualize CPU consumption of profiled code-paths of software.
They are useful visualization tools in many types of performance
investigations.

Luca Canali has previously shown the benefits of using flamegraph
visualizations of Oracle query plans for debugging slow database queries.
Pg_flame is in extension of that work for Postgres.

This tool relies on the
[`d3-flame-graph`](https://github.com/spiermar/d3-flame-graph) plugin to
generate the flamegraph.
