#!/usr/bin/env bash

go run cmd/lector/main.go 1 ms/users.txt &
go run cmd/lector/main.go 2 ms/users.txt &
go run cmd/escritor/main.go 3 ms/users.txt &
go run cmd/escritor/main.go 4 ms/users.txt &
