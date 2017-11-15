# disco-http

[![Go Report Card](https://goreportcard.com/badge/github.com/continuul/disco-http)](https://goreportcard.com/report/github.com/continuul/disco-http)

A simple demonstration of Hashicorp Serf and Raft libraries to illustrate
a peer-to-peer Raft-based network of nodes that self-heal, and perform
leadership election.

This is a reference implementation of the [Hashicorp Raft implementation v1.0](https://github.com/hashicorp/raft).
[Raft](https://raft.github.io/) is a _distributed consensus protocol_, meaning
its purpose is to ensure that a set of nodes -- a cluster -- agree on the
state of some arbitrary state machine, even when nodes are vulnerable to
failure and network partitions. Distributed consensus is a fundamental
concept when it comes to building fault-tolerant systems.

A simple example system like disco-http makes it easy to study the Raft
consensus protocol in general, and Hashicorp's Raft implementation in particular.

## Reading and writing keys

The reference implementation is a very simple in-memory key-value store.
You can set a key by sending a request to the HTTP bind address
(which defaults to `localhost:11000`):

```bash
curl -XPOST localhost:11000/key -d '{"foo": "bar"}'
```

You can read the value for a key like so:
```bash
curl -XGET localhost:11000/key/foo
```

## Building disco-http

Starting and running a disco-http cluster is easy. Download disco-http like so:

```bash
mkdir demos
cd demos/
export GOPATH=$PWD
go get continuul.io/disco-http
```

## Running disco-http

Building disco-http requires Go 1.9 or later.

Run your first disco-http node like so:
```bash
disco-http -id node0 /tmp/node0
```

You can now set a key and read its value back:
```bash
curl -XPOST localhost:11000/key -d '{"user1": "batman"}'
curl -XGET localhost:11000/key/user1
```

### Bring up a cluster

_A walkthrough of setting up a more realistic cluster is [here](CLUSTERING.md)._

Let's bring up 2 more nodes, so we have a 3-node cluster. That way we can tolerate the failure of 1 node:
```bash
disco-http -id node1 -client :11001 -bind :12001 -join :11000 /tmp/node1
disco-http -id node2 -client :11002 -bind :12002 -join :11000 /tmp/node2
```
_This example shows each disco-http node running on the same host, so each node must listen on different ports. This would not be necessary if each node ran on a different host._

This tells each new node to join the existing node. Once joined, each node now knows
about the key:
```bash
curl -XGET localhost:11000/key/user1
curl -XGET localhost:11001/key/user1
curl -XGET localhost:11002/key/user1
```

Furthermore you can add a second key:
```bash
curl -XPOST localhost:11000/key -d '{"user2": "robin"}'
```

Confirm that the new key has been set like so:
```bash
curl -XGET localhost:11000/key/user2
curl -XGET localhost:11001/key/user2
curl -XGET localhost:11002/key/user2
```

#### Stale reads

Because any node will answer a GET request, and nodes may "fall behind"
updates, stale reads are possible. Again, disco-http is a simple program,
for the purpose of demonstrating a distributed key-value store. If you
are particularly interested in learning more about issue, you should
check out [rqlite](https://github.com/rqlite/rqlite). rqlite allows the
client to control [read consistency](https://github.com/rqlite/rqlite/blob/master/doc/CONSISTENCY.md),
allowing the client to trade off read-responsiveness and correctness.

Read-consistency support could be ported to disco-http if necessary.

### Tolerating failure

Kill the leader process and watch one of the other nodes be elected leader.
The keys are still available for query on the other nodes, and you can set
keys on the new leader. Furthermore, when the first node is restarted, it
will rejoin the cluster and learn about any updates that occurred while it
was down.

A 3-node cluster can tolerate the failure of a single node, but a 5-node
cluster can tolerate the failure of two nodes. But 5-node clusters require
that the leader contact a larger number of nodes before any change e.g. setting
a key's value, can be considered committed.

### Leader-forwarding

Automatically forwarding requests to set keys to the current leader is not
implemented. The client must always send requests to change a key to the
leader or an error will be returned.
