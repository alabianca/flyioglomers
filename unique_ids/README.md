# Challenge 2: Unique ID Generation

Challenge page: <https://fly.io/dist-sys/2/>

## Problem

Every node must answer `generate` requests with a globally unique ID. The
service has to stay totally available, meaning it keeps answering even while
the network is partitioned and nodes cannot talk to each other. That rules out
any scheme where nodes coordinate to hand out IDs.

## Approach

The node uses a [Snowflake](https://en.wikipedia.org/wiki/Snowflake_ID)-style
generator, implemented in [`internal/snowflake`](../internal/snowflake/). Each
ID is a 64-bit integer with this layout:

| Bits | Field | Notes |
|------|-------|-------|
| 41 | Timestamp | Milliseconds since Twitter's Snowflake epoch (2010-11-04) |
| 10 | Node ID | Parsed from the Maelstrom node name, so `n3` becomes `3` |
| 12 | Sequence | Counter for IDs generated within the same millisecond |

Uniqueness falls out of the layout rather than any cross-node agreement. Two
nodes can never collide because their node ID bits differ, and a single node
never collides with itself because the sequence counter separates IDs minted
in the same millisecond. Partitions are irrelevant.

A few details worth noting:

- **Sequence overflow.** If a node mints more than 4096 (largest 12 bit integer) IDs in one millisecond
  the generator borrows the next millisecond and resets the sequence, instead
  of blocking or failing.
- **Clock going backwards.** If the wall clock reads earlier than the last
  timestamp used, the generator keeps using the last timestamp so IDs stay
  monotonic within a node.
- **Thread safety.** A mutex guards the counter and last-timestamp state, since
  the Maelstrom library runs handlers concurrently.
- **Node ID range.** The 10-bit field supports node IDs 0 through 1023. Larger
  IDs return an error rather than silently wrapping.

The generator has unit tests covering timestamp and node ID round-tripping,
sequence increments within a millisecond, sequence reset on a new millisecond,
and the overflow case.

## Gotchas

### Send IDs as strings, not JSON integers

The first version of the node replied with the raw 64-bit integer in `body.id`.
That is legal as far as the workload is concerned. Maelstrom only checks that
IDs are unique, and it accepts integers, floats, or strings for the value.

The problem is the size of the number. A Snowflake ID uses the full 64-bit
range, which is far beyond the 53 bits of integer precision that many JSON
parsers assume. Once the IDs crossed that threshold Maelstrom started to see
problems with the integer values in its JSON handling, and the run failed even
though the generator itself was producing unique IDs.

The fix was to stringify the ID before putting it in the reply:

```go
body["id"] = id.String()
```

Strings pass through JSON untouched, so the exact 64-bit value survives the
round trip and the uniqueness check works as intended. The lesson generalises
beyond this challenge: when an ID or counter can exceed 2^53, send it as a
string over JSON.

## Running

```sh
./run unique_ids
go test ./internal/snowflake/
```

Workload flags: `-w unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition`
