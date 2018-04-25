# SCP Nomination Prototype

Prototype implementation of nomination protocol in [stellar-core](https://github.com/stellar/stellar-core)

## Installation

To install and deploy the source, you need to install these packages,

 - golang

Clone this repository and run.

## Test Run

There are 
```sh
$ cd src/ledger
$ go test
```

## Deployment

```sh
$ cd src/run
$ go run run-ledger.go -h
Usage of /var/folders/9w/lx5fhzd54z916pmyt2hh2jsr0000gn/T/go-build528634220/b001/exe/run-ledger:
  -conf string
    	input yml conf file (default "../config/single_quorum.yml")
```

### Configuration File

Set the config file.
```json
default:
    threshold: 5

node:
    n0:
        kind: 0         # node kind. 0: normal node, 1: faulty node
        name: n0        # node name
        validators:     # node validators
            - n1
            - n2
            - n3
            - n4
        messages:       # messages in node message pool
            - message from n0
            - message2 from n0
    n1:
        kind: 0 # normal
        name: n1
        validators:
            - n0
            - n2
            - n3
            - n4
    n2:
        kind: 0 # normal
        name: n2
        validators:
            - n0
            - n1
            - n3
            - n4
    n3:
        kind: 0 # normal
        name: n3
        validators:
            - n0
            - n1
            - n2
            - n4
    n4:
        kind: 0 # normal
        name: n4
        validators:
            - n0
            - n1
            - n2
            - n3
        messages:
            - message3 from n4
            - message4 from n4
```

### Running `run-ledger.go`

```sh
$ go run run-ledger.go -conf ../config/single_quorum.yml
```
