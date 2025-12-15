## Purpose
Client - database - server setup to generate traffic in a work load.

## Requirements
Maria DB (reasonably recent is fine)
Go version 1.21 or above

## How to run
(recommend to use tmux)

# get the server going first
terminal 1 > `go run server.go`

# client calls every N seconds
terminal 2 > `while true; do go run client.go; sleep 2; done`

