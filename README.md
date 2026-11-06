# Fly.io Gossip Glomers

A learning project working through the [Gossip Glomers](https://fly.io/dist-sys/)
distributed systems challenges from Fly.io, written in Go.

Each challenge is a small node program that speaks the
[Maelstrom](https://github.com/jepsen-io/maelstrom) protocol over stdin/stdout.
Maelstrom runs a cluster of these nodes, sends them requests, injects network
partitions and latency, and checks the results with
[Jepsen](https://github.com/jepsen-io/jepsen). The point of the project is to
build intuition for the hard parts of distributed systems by getting them wrong
first and then fixing them.

## Progress

| # | Challenge | Status | Directory |
|---|-----------|--------|-----------|
| 1 | [Echo](https://fly.io/dist-sys/1/) | ✅ Done | [`echo/`](echo/) |
| 2 | [Unique ID Generation](https://fly.io/dist-sys/2/) | ✅ Done | [`unique_ids/`](unique_ids/) |
| 3a | [Single-Node Broadcast](https://fly.io/dist-sys/3a/) | ⬜ Not started | |
| 3b | [Multi-Node Broadcast](https://fly.io/dist-sys/3b/) | ⬜ Not started | |
| 3c | [Fault Tolerant Broadcast](https://fly.io/dist-sys/3c/) | ⬜ Not started | |
| 3d | [Efficient Broadcast, Part I](https://fly.io/dist-sys/3d/) | ⬜ Not started | |
| 3e | [Efficient Broadcast, Part II](https://fly.io/dist-sys/3e/) | ⬜ Not started | |
| 4 | [Grow-Only Counter](https://fly.io/dist-sys/4/) | ⬜ Not started | |
| 5a | [Single-Node Kafka-Style Log](https://fly.io/dist-sys/5a/) | ⬜ Not started | |
| 5b | [Multi-Node Kafka-Style Log](https://fly.io/dist-sys/5b/) | ⬜ Not started | |
| 5c | [Efficient Kafka-Style Log](https://fly.io/dist-sys/5c/) | ⬜ Not started | |
| 6a | [Single-Node, Totally-Available Transactions](https://fly.io/dist-sys/6a/) | ⬜ Not started | |
| 6b | [Totally-Available, Read Uncommitted Transactions](https://fly.io/dist-sys/6b/) | ⬜ Not started | |
| 6c | [Totally-Available, Read Committed Transactions](https://fly.io/dist-sys/6c/) | ⬜ Not started | |

## Challenge notes

### 1. Echo

The "hello world" of Maelstrom. The node receives an `echo` message and replies
with an `echo_ok` carrying the same payload. Mostly a check that the toolchain,
the Go Maelstrom library, and the test harness all work together.

### 2. Unique ID Generation

Every node must hand out globally unique IDs with no coordination, even while
the network is partitioned. The solution is a
[Snowflake](https://en.wikipedia.org/wiki/Snowflake_ID)-style generator in
`internal/snowflake`: each 64-bit ID packs a millisecond timestamp,
the node's ID, and a per-millisecond sequence counter. Uniqueness falls out of
the layout rather than any cross-node agreement, so partitions are irrelevant.
The generator has its own unit tests covering sequence rollover and timestamp
changes.

## Layout

```
.
├── go.mod         # Single Go module for the whole repo (module name: glomers)
├── run            # Build a challenge and run its Maelstrom workload
├── echo/          # Challenge 1
├── unique_ids/    # Challenge 2
├── internal/      # Packages shared across challenges (e.g. internal/snowflake)
├── bin/           # Built node binaries (gitignored)
├── maelstrom/     # Maelstrom release, downloaded locally (gitignored)
└── store/         # Maelstrom test output, one directory per run (gitignored)
```

Each challenge directory is a `package main` that builds a single node binary.
Shared code lives under `internal/` and is imported as `glomers/internal/...`.
The module depends on `github.com/jepsen-io/maelstrom/demo/go`, Fly's thin Go
wrapper around the Maelstrom protocol.

## Running the tests

### Prerequisites

- Go 1.22 or newer
- A JVM (OpenJDK 17+) for Maelstrom
- Graphviz and gnuplot, if you want Maelstrom's plots

On macOS:

```sh
brew install go openjdk graphviz gnuplot
```

### Install Maelstrom

The Maelstrom release is not checked in. Download the latest tarball from the
[releases page](https://github.com/jepsen-io/maelstrom/releases) and unpack it
into `maelstrom/` at the root of this repo:

```sh
curl -LO https://github.com/jepsen-io/maelstrom/releases/download/v0.2.3/maelstrom.tar.bz2
tar -xjf maelstrom.tar.bz2
rm maelstrom.tar.bz2
```

That leaves `./maelstrom/maelstrom` as the runner used in the commands below.

### Build and test a challenge

The `run` script builds a challenge into `bin/<challenge>` and then runs the
matching Maelstrom workload with the flags from the Fly.io challenge page:

```sh
./run echo
./run unique_ids
```

Any extra arguments are passed straight to Maelstrom, so you can override or
extend the defaults:

```sh
./run unique_ids --time-limit 60
```

To add a new challenge, create a `<challenge>/main.go` and add a case with its
workload flags to the `run` script. The equivalent manual invocation looks like:

```sh
go build -o bin/unique_ids ./unique_ids
./maelstrom/maelstrom test -w unique-ids --bin bin/unique_ids \
  --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
```

### Inspecting results

Maelstrom writes every run to `store/<workload>/<timestamp>/`, with
`store/latest` pointing at the most recent one. The quickest way to dig in is the
built-in web UI:

```sh
./maelstrom/maelstrom serve
```

Then open <http://localhost:8080> to browse timelines, latency graphs, and the
raw message log for each run.

### Unit tests

Shared packages under `internal/` have plain Go tests alongside the Maelstrom
runs:

```sh
go test ./...
```

## References

- [Gossip Glomers](https://fly.io/dist-sys/) — the challenge series
- [Maelstrom](https://github.com/jepsen-io/maelstrom) — the test workbench
- [Maelstrom protocol](https://github.com/jepsen-io/maelstrom/blob/main/doc/protocol.md) and
  [workloads](https://github.com/jepsen-io/maelstrom/blob/main/doc/workloads.md) — the docs
- [Jepsen](https://jepsen.io/) — the consistency checker underneath it all
