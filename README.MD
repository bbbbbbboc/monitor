# Flow Assignment

## Feature Request
Build a cli supports monitoring transactions related to a given account.

## How to Run
1. If you don't see the `monitor` executable, first build it by running `go install` and `go build`
1. `monitor` and `monitor transactions` will list help
2. `monitor transactions <address>` will list all tx of the given address from now on.

## Tech Stack
This cli is built using go, cli framework [cobra](https://cobra.dev/) and [flow go sdk](https://docs.onflow.org/flow-go-sdk/).

## Development
1. Put the repo to your GOPATH and create a directory for it, for me I put it in hone/go/src/flow/monitor
1. `go get` to install dependency, this will install cobra and flow sdk.
2. `go install` and `go build` to generate the executable `monitor`

## Test
Some test addresses to use,
0xead892083b3e2c6c,
0xaf141341c08b59e0,
0x05b75c9eb5c80ce0,
0x55ad22f01ef568a1,
0x7b0aac240b2bacda

I'm using 0x55ad22f01ef568a1 since it's pretty active
after run monitor, check [block explorer](https://flowscan.org/) to see if every transaction shows in cli.
