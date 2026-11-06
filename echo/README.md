# Challenge 1: Echo

Challenge page: <https://fly.io/dist-sys/1/>

## Problem

The node receives an `echo` message and must reply with an `echo_ok` message
carrying the same `echo` payload back to the client.

```json
{ "type": "echo", "msg_id": 1, "echo": "Please echo 35" }
{ "type": "echo_ok", "msg_id": 1, "in_reply_to": 1, "echo": "Please echo 35" }
```

## Approach

This is the "hello world" of Maelstrom. The handler unmarshals the body into a
loosely typed map, flips `type` to `echo_ok`, and replies with the same body.
The Go Maelstrom library takes care of `msg_id` and `in_reply_to`.

The real purpose of this challenge is to check that the toolchain works end to
end: Go builds the node, Maelstrom launches it, messages flow over stdin and
stdout, and the Jepsen checker reads the results.

## Running

```sh
./run echo
```

Workload flags: `-w echo --node-count 1 --time-limit 10`
