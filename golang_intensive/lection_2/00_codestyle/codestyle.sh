#!/bin/bash

# go get golang.org/x/tools/cmd/goimports

gofmt -s -w main.go
goimports -w main.go
