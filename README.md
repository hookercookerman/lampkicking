# Overview

LampKicking Test

Ok there are 3 main packages here.

* Persistar (persistence abstraction service)
* Persist (client for Persistar)
* User (user service)

# Persistence Abstraction (Persistar)

* Constrained to be a JSON document store.
* Changeable Durability (Redis Default) Composable Adaptor Interface for keyvalue
  and graph.
* Graph api for relation data abstractions, well that was the idea
* Use Notion of services instead of handlers
* Interactors are like use case services, I just like to abstract to these when
  building apps

# Persist

* Simple Client lib for Persistar (style taken from stripe-go).

# User Service

* Leverages Persist and Persistar. (which means should be stupid simple)
* Provides the 3 endpoints required for the test.

# Getting Started

Using GoDeps

https://github.com/tools/godep


Redis Server will need to be running

if you want to quickly run both service

```shell
foreman start
```
* create some users

```shell
curl localhost:9001/users -X POST -d '{"name" : "egg"}'
curl localhost:9001/users -X POST -d '{"name" : "beans"}'
curl localhost:9001/users -X POST -d '{"name" : "cheese"}'
```

response

```json
{"id" : "id"}
```

* Add Connections
```shell
curl localhost:9001/users/:user_id/connections/:other_user_id -X PUT
```

* Get connections
```shell
curl http://localhost:9001/users/:user_id/connections
```

# Tests

sorry Ginkgo I have been rspec for a long time

```shell
ginkgo -r
```

These could be alot better/more especially the persist client which has none lol; just
run out of time,

the graph service should not talk redis in terms of data setup etc but should be
using adaptor directly or it should be mocked; so thats not great

# Improvements

* Investigate http://www.grpc.io/ better then json apis!
* Better http code handlers aka 422 instead of 500 all the time; and probably
  shouldnt be 500 in first place
* More concrete error types generally around the durability in the
  persistar
