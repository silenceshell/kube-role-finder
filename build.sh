#!/bin/bash

# for go before 1.11
#GOPATH=/home/bottle/Code/Go/ GO111MODULE=off go build

# for go after 1.11
GOPROXY="https://gocenter.io" GO111MODULE=on go build
