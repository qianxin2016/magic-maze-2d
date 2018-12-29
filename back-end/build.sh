#!/bin/bash
export GOPATH=`pwd`
go build -o magic-maze src/main.go src/maze.go src/mysql.go
