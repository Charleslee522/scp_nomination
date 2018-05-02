# SCP Nomination Prototype

Prototype implementation of nomination protocol in [stellar-core](https://github.com/stellar/stellar-core)

## Installation(Mac OS X)

To install and deploy the source, you need to install these packages,

 - go: 1.10.1 or higher

 ```
 $ brew update
 $ brew install golang

 $ export GOPATH=$HOME/go # set your path
 $ export GOROOT=/usr/local/opt/go/libexec
 $ export PATH=$PATH:$GOPATH/bin
 $ export PATH=$PATH:$GOROOT/bin
 ```

Clone this repository and run.

## Clone Repository

```
go get github.com/Charleslee522/scp_nomination
```

## Test Run

```sh
$ cd $GOPATH/src/github.com/Charleslee522/scp_nomination
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
